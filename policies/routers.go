package policies

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

func PolicyRegister(router *gin.RouterGroup) {
	router.GET("", PolicyList)
	router.GET("/:name", PolicyRetrieve)
	router.POST("/:name", PolicyCreateOrUpdate)
	router.PUT("/:name", PolicyCreateOrUpdate)
	router.DELETE("/:name", PolicyDelete)
}

func PolicyList(c *gin.Context) {
	l, _ := common.GetLogger()
	list, err := strconv.ParseBool(c.Query("list"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("policy", err))
		return
	}
	if !list {
		c.JSON(http.StatusMethodNotAllowed, common.NewError("policy", errors.New("method not allowed")))
		return
	}
	policyModels, count, err := FindManyPolicies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	l.Debug("Retrieved models", "count", count, "models", policyModels, "err", err)
	serializer := PoliciesSerializer{C: c, Policies: policyModels}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func PolicyRetrieve(c *gin.Context) {
	name := c.Param("name")
	policyModel, err := FindOnePolicy(&PolicyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("policy", errors.New("Invalid policy name")))
		return
	}
	policySerializer := PolicySerializer{C: c, PolicyModel: policyModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, policySerializer.Response()))
}

func PolicyCreateOrUpdate(c *gin.Context) {
	name := c.Param("name")
	policyModel, err := FindOnePolicy(&PolicyModel{Name: name})
	if err != nil {
		policyModel = PolicyModel{Name: name}
	}
	policyModelValidator := NewPolicyModelValidatorFillWith(policyModel)
	if err := policyModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("policy-validator", err))
		return
	}

	if policyModel.ID != 0 {
		policyModelValidator.policyModel.ID = policyModel.ID
		if err := policyModel.Update(policyModelValidator.policyModel); err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
	} else {
		if err := SaveOne(&policyModelValidator.policyModel); err != nil {
			c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}
	}
	policySerializer := PolicySerializer{C: c, PolicyModel: policyModelValidator.policyModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, policySerializer.Response()))
}

func PolicyDelete(c *gin.Context) {
	name := c.Param("name")
	err := DeletePolicyModel(&PolicyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("policy", errors.New("Invalid policy name")))
		return
	}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, common.NewStatusResponse(http.StatusOK, "ok")))
}
