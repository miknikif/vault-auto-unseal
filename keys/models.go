package keys

import (
	"github.com/jinzhu/gorm"
	"github.com/miknikif/vault-auto-unseal/common"
)

type KeyType string

const (
	KEY_TYPE_AES256_GCM96 KeyType = "aes256-gcm96"
)

type AESKeyModel struct {
	gorm.Model
	KeyID   uint
	Name    string `gorm:"uniqueIndex:keynamever;"`
	Version int    `gorm:"uniqueIndex:keynamever;"`
	AESKey  string
}

type KeyModel struct {
	gorm.Model
	KeyID                uint
	Name                 string
	Type                 KeyType
	Keys                 []AESKeyModel `gorm:"foreignKey:KeyID"`
	AllowPlaintextBackup bool
	AutoRotatePeriod     int
	DeletionAllowed      bool
	Derived              bool
	Exportable           bool
	ImportedKey          bool
	LatestVersion        int
	MinAvailableVersion  int
	MinDecryptionVersion int
	MinEncryptionVersion int
	SupportsDecryption   bool
	SupportsDerivation   bool
	SupportsEncryption   bool
	SupportsSigning      bool
}

func (model *KeyModel) Update(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Updating KeyModel")
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Model(model).Update(data).Error
	return err
}

func SaveOne(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Saving data to DB")
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Save(data).Error
	return err
}

func FindOneAESKey(condition interface{}) (AESKeyModel, error) {
	var model AESKeyModel
	l, err := common.GetLogger()
	if err != nil {
		return model, err
	}
	l.Debug("Starting retrieval of the AESKeyModel from the DB")
	db, err := common.GetDB()
	if err != nil {
		return model, err
	}
	tx := db.Begin()
	tx.Where(condition).First(&model)
	err = tx.Commit().Error
	l.Debug("Finished retrieval of the AESKeyModel from the DB")
	return model, err
}

// func FindManyAESKey(name string, limit int, offset int) ([]AESKeyModel, int, error) {
// 	var models []AESKeyModel
//     var model AESKeyModel
// 	var count int
// 	if name == "" {
// 		return models, count, errors.New("Key name should be specified, but received empty name")
// 	}
// 	c, err := common.GetConfig()
// 	if err != nil {
// 		return models, count, err
// 	}
// 	c.Logger.Debug("Starting retrieval of the AESKeyModels from the DB")
// 	db := c.DB
// 	tx := db.Begin()
//
//     tx.Where(AESKeyModel{Name: name}).First(&model)
//
//     if model.KeyID != 0 {
//         count = tx.Model(&model)
//     }
//
// 	tx.Where(condition).First(&model)
// 	err = tx.Commit().Error
// 	c.Logger.Debug("Finished retrieval of the AESKeyModels from the DB")
// 	return model, err
// }

func FindOneKey(condition interface{}) (KeyModel, error) {
	var model KeyModel
	l, err := common.GetLogger()
	if err != nil {
		return model, err
	}
	l.Debug("Starting retrieval of the KeyModel from the DB")
	db, err := common.GetDB()
	if err != nil {
		return model, err
	}
	tx := db.Begin()
	tx.Where(condition).First(&model)
	tx.Model(&model).Related(&model.Keys, "Keys")
	err = tx.Commit().Error
	l.Debug("Finished retrieval of the KeyModel from the DB")
	return model, err
}

func DeleteAESKeyModel(condition interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Deleting the AESKeyModel from the DB")
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Where(condition).Delete(AESKeyModel{}).Error
	return err
}

func DeleteKeyModel(condition interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Deleting the KeyModel from the DB")
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Where(condition).Delete(KeyModel{}).Error
	return err
}
