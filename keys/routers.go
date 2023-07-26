package keys

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"net/http"
)

func KeysRegister(router *gin.RouterGroup) {
	router.GET("/", KeyList)
	router.GET("/:slug", KeyRetrieve)
	router.POST("/:slug", KeyCreate)
	router.PUT("/:slug", KeyUpdate)
	router.DELETE("/:slug", KeyDelete)
}

func KeyCreate(c *gin.Context) {
	keyModelValidator := NewKeyModelValidator()
	if err := keyModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("keys", err))
		return
	}

	if err := SaveOne(&keyModelValidator.keyModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	serializer := KeySerializer{c, keyModelValidator.keyModel}
	c.JSON(http.StatusCreated, gin.H{"key": serializer.Response()})
}

func KeyList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, common.NewError("keys", errors.New("Path not implemented yet.")))
}

func KeyRetrieve(c *gin.Context) {
	slug := c.Param("slug")
	keyModel, err := FindOneKey(&KeyModel{Name: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}
	serializer := KeySerializer{c, keyModel}
	c.JSON(http.StatusOK, gin.H{"key": serializer.Response()})
}

func KeyUpdate(c *gin.Context) {
	slug := c.Param("slug")
	keyModel, err := FindOneKey(&KeyModel{Name: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}
	keyModelValidator := NewKeyModelValidatorFillWith(keyModel)
	if err := keyModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("keys", err))
		return
	}

	keyModelValidator.keyModel.KeyID = keyModel.KeyID
	if err := keyModel.Update(keyModelValidator); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := KeySerializer{c, keyModel}
	c.JSON(http.StatusOK, gin.H{"key": serializer.Response()})
}

func KeyDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := DeleteKeyModel(&KeyModel{Name: slug})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": "Delete success"})
}
