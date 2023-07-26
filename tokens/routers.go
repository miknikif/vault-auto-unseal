package tokens

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

func TokenRegister(router *gin.RouterGroup) {
	router.POST("/create", TokenCreate)
	router.POST("/lookup", TokenRetrieve)
	router.POST("/lookup-accessor", TokenRetrieve)
}

// {"ttl":"0s","explicit_max_ttl":"0s","period":"0s","display_name":"","num_uses":0,"renewable":true,"type":"service","entity_alias":""}
// {"policies":["default"],"ttl":"0s","explicit_max_ttl":"0s","period":"24h0m0s","display_name":"","num_uses":0,"renewable":true,"type":"service","entity_alias":""}
func TokenCreate(c *gin.Context) {
	l, _ := common.GetLogger()
	tokenModelValidator := NewTokenModelValidator()
	if err := tokenModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("tokens", err))
		return
	}
	l.Debug("Saving token to the DB: ", "token", tokenModelValidator.tokenModel)
	if err := SaveOne(&tokenModelValidator.tokenModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := TokenSerializer{C: c, TokenModel: tokenModelValidator.tokenModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func TokenRetrieve(c *gin.Context) {
	l, _ := common.GetLogger()
	tokenLookupModelValidator := NewTokenLookupModelValidator()
	if err := tokenLookupModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("tokens", err))
		return
	}
	tokenModel, err := FindOneToken(&tokenLookupModelValidator.tokenModel)
	l.Debug("Retrieved token:", "token", tokenModel, "err", err)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("tokens", errors.New("Specified token not found")))
		return
	}
	if tokenLookupModelValidator.findWithAccessor {
		tokenModel.TokenID = ""
	}
	serializer := TokenSerializer{C: c, TokenModel: tokenModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}
