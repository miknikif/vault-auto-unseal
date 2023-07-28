package keys

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

func KeysOperationsRegister(router *gin.RouterGroup) {
	router.PUT("/encrypt/:name", EncryptData)
	router.PUT("/decrypt/:name", DecryptData)
	router.PUT("/rewrap/:name", RewrapData)
	KeysRegister(router.Group("/keys"))
}

func KeysRegister(router *gin.RouterGroup) {
	router.GET("", KeyList)
	router.GET("/:name", KeyRetrieve)
	router.PUT("/:name", KeyCreate)
	router.DELETE("/:name", KeyDelete)
	router.PUT("/:name/config", KeyUpdate)
	router.PUT("/:name/rotate", KeyRotate)
}

func KeyCreate(c *gin.Context) {
	if allowed := common.VerifyCreateAccess(c); !allowed {
		c.JSON(http.StatusForbidden, common.NewError("auth", errors.New("permission denied")))
		return
	}
	name := c.Param("name")
	keyModelValidator := NewKeyModelValidator()
	keyModelValidator.Name = name
	if err := keyModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("keys", err))
		return
	}

	if keyModel, err := FindOneKey(&KeyModel{Name: name}); err == nil || keyModel.ID != 0 {
		c.JSON(http.StatusConflict, common.NewError("keys", errors.New("key already exist")))
		return
	}

	if err := SaveOne(&keyModelValidator.keyModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	serializer := KeySerializer{c, keyModelValidator.keyModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func KeyList(c *gin.Context) {
	if allowed := common.VerifyListAccess(c); !allowed {
		c.JSON(http.StatusForbidden, common.NewError("auth", errors.New("permission denied")))
		return
	}
	l, _ := common.GetLogger()
	list := common.ParseBool(c.Query("list"), false)
	if !list {
		c.JSON(http.StatusMethodNotAllowed, common.NewError("keys", errors.New("method not allowed")))
		return
	}
	keyModels, count, err := FindManyKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}
	l.Debug("Retrieved models", "count", count, "models", keyModels, "err", err)
	serializer := KeysSerializer{C: c, Keys: keyModels}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func KeyRetrieve(c *gin.Context) {
	name := c.Param("name")
	keyModel, err := FindOneKey(&KeyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}
	serializer := KeySerializer{c, keyModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func KeyUpdate(c *gin.Context) {
	name := c.Param("name")
	keyModel, err := FindOneKey(&KeyModel{Name: name})
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
	if err := keyModel.Update(keyModelValidator.keyModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := KeySerializer{c, keyModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func KeyRotate(c *gin.Context) {
	name := c.Param("name")
	keyModel, err := FindOneKey(&KeyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}
	key, err := createNewAESKeyModel(keyModel.LatestVersion+1, AES_KEY_SIZE_256)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("keys", err))
		return
	}
	keyModel.Keys = append(keyModel.Keys, key)
	keyModelValidator := NewKeyModelValidatorFillWith(keyModel)
	if err := keyModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("keys", err))
		return
	}

	keyModelValidator.keyModel.KeyID = keyModel.KeyID
	if err := keyModel.Update(keyModelValidator.keyModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := KeySerializer{c, keyModel}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func KeyDelete(c *gin.Context) {
	name := c.Param("name")
	keyModel, err := FindOneKey(&KeyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}

	if !keyModel.DeletionAllowed {
		c.JSON(http.StatusForbidden, common.NewError("keys", errors.New("key is marked as not deletable")))
		return
	}

	err = DeleteKeyModel(&keyModel)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("database", errors.New("Unable to delete key")))
		return
	}
	c.JSON(http.StatusOK, common.NewStatusResponse(http.StatusOK, "key deleted"))
}

func EncryptData(c *gin.Context) {
	name := c.Param("name")
	encryptDataValidator := NewEncryptDataValidator()
	if err := encryptDataValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("encrypt", err))
		return
	}
	keyModel, err := FindOneKey(&KeyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("Key not found")))
		return
	}

	key, err := findKeyVersion(keyModel.Keys, keyModel.LatestVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", fmt.Errorf("key version v%d doesn't exist", keyModel.LatestVersion)))
		return
	}

	if err := encryptDataWithAES(key, &encryptDataValidator.aesPayload); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", err))
		return
	}

	serializer := EncryptDataSerializer{C: c, AESPayload: encryptDataValidator.aesPayload}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func DecryptData(c *gin.Context) {
	name := c.Param("name")
	decryptDataValidator := NewDecryptDataValidator()
	if err := decryptDataValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("decrypt", err))
		return
	}
	keyModel, err := FindOneKey(&KeyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("key not found")))
		return
	}

	if decryptDataValidator.aesPayload.Version < keyModel.MinDecryptionVersion {
		c.JSON(http.StatusForbidden, common.NewError("keys", fmt.Errorf("minimum version to decrypt is v%d, but you're requested v%d to be decrypted", keyModel.MinDecryptionVersion, decryptDataValidator.aesPayload.Version)))
		return
	}

	key, err := findKeyVersion(keyModel.Keys, decryptDataValidator.aesPayload.Version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", fmt.Errorf("key version v%d doesn't exist", decryptDataValidator.aesPayload.Version)))
		return
	}

	if err := decryptDataWithAES(key, &decryptDataValidator.aesPayload); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", err))
		return
	}

	serializer := DecryptDataSerializer{C: c, AESPayload: decryptDataValidator.aesPayload}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}

func RewrapData(c *gin.Context) {
	name := c.Param("name")
	decryptDataValidator := NewDecryptDataValidator()
	if err := decryptDataValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("decrypt", err))
		return
	}
	keyModel, err := FindOneKey(&KeyModel{Name: name})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("keys", errors.New("key not found")))
		return
	}

	if decryptDataValidator.aesPayload.Version < keyModel.MinDecryptionVersion {
		c.JSON(http.StatusForbidden, common.NewError("keys", fmt.Errorf("minimum version to decrypt is v%d, but you're requested v%d to be decrypted", keyModel.MinDecryptionVersion, decryptDataValidator.aesPayload.Version)))
		return
	}

	key, err := findKeyVersion(keyModel.Keys, decryptDataValidator.aesPayload.Version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", fmt.Errorf("key version v%d doesn't exist", decryptDataValidator.aesPayload.Version)))
		return
	}

	if err := decryptDataWithAES(key, &decryptDataValidator.aesPayload); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", err))
		return
	}

	encryptDataValidator := NewEncryptDataValidatorFillWith(decryptDataValidator.aesPayload)
	if err := encryptDataValidator.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("encrypt", err))
		return
	}

	keyLatest, err := findKeyVersion(keyModel.Keys, keyModel.LatestVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", fmt.Errorf("key version v%d doesn't exist", keyModel.LatestVersion)))
		return
	}

	if err := encryptDataWithAES(keyLatest, &encryptDataValidator.aesPayload); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("transit", err))
		return
	}

	serializer := EncryptDataSerializer{C: c, AESPayload: encryptDataValidator.aesPayload}
	c.JSON(http.StatusOK, common.NewGenericResponse(c, serializer.Response()))
}
