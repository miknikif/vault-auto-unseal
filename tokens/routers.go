package tokens

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

func TokenRegister(router *gin.RouterGroup) {
	router.POST("/create", TokenCreate)
	router.POST("/renew", TokenRenew)
	router.PUT("/renew", TokenRenew)
	router.POST("/renew-accessor", TokenRenew)
	router.GET("/lookup-self", TokenSelfRetrieve)
	router.POST("/lookup", TokenRetrieve)
	router.POST("/lookup-accessor", TokenRetrieve)
	router.POST("/revoke", TokenDelete)
	router.PUT("/revoke", TokenDelete)
	router.POST("/revoke-accessor", TokenDelete)
}

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

func TokenRenew(c *gin.Context) {
	l, _ := common.GetLogger()
	tokenLookupModelValidator := NewTokenLookupModelValidator(false)
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
	if err := tokenModel.Renew(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("tokens", err))
		return
	}
	if tokenLookupModelValidator.findWithAccessor {
		tokenModel.TokenID = ""
	}
	serializer := TokenSerializer{C: c, TokenModel: tokenModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func TokenSelfRetrieve(c *gin.Context) {
	l, _ := common.GetLogger()
	tokenLookupModelValidator := NewTokenLookupModelValidator(true)
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

func TokenRetrieve(c *gin.Context) {
	l, _ := common.GetLogger()
	tokenLookupModelValidator := NewTokenLookupModelValidator(false)
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

func TokenDelete(c *gin.Context) {
	l, _ := common.GetLogger()
	tokenLookupModelValidator := NewTokenLookupModelValidator(false)
	if err := tokenLookupModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("tokens", err))
		return
	}
	tokenModel, err := FindOneToken(&tokenLookupModelValidator.tokenModel)
	l.Debug("TokenDelete: found existing token", "token", tokenModel)
	if err != nil {
		c.JSON(http.StatusOK, common.NewGenericResponse(c, nil))
		return
	}
	if err := DeleteTokenModel(&tokenModel); err != nil {
		l.Debug("TokenDelete: unable to delete the token", "token", tokenModel, "err", err)
		c.JSON(http.StatusOK, common.NewGenericResponse(c, nil))
		return
	}
	l.Debug("TokenDelete: token deleted")
	c.JSON(http.StatusNoContent, nil)
}
