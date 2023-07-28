package keys

import (
	"errors"
	"fmt"

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
	Name    int
	Version int `gorm:"uniqueIndex:keynamever;"`
	AESKey  AESKey
}

type KeyModel struct {
	gorm.Model
	KeyID                uint
	Name                 string
	Type                 KeyType
	Keys                 []AESKeyModel `gorm:"foreignKey:KeyID;constraint:OnDelete:CASCADE;"`
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

type AESPayload struct {
	Plaintext string
	Pref      string
	Version   int
	Payload   string
}

func (s *AESPayload) validatePlaintext() error {
	l, _ := common.GetLogger()
	l.Debug("AESPayload.validatePlaintext - started", "self", s)
	if s.Plaintext == "" {
		return errors.New("plaintext is empty")
	}
	if _, err := common.DecFromB64(s.Plaintext); err != nil {
		return errors.New("plaintext should be b64 encoded")
	}
	l.Debug("AESPayload.validatePlaintext - ended", "self", s)
	return nil
}

func (s *AESPayload) validateCiphertext() error {
	if s.Pref == "" {
		return errors.New("ciphertext prefix is empty")
	}
	if s.Version == 0 {
		return errors.New("wrong ciphertext version specified")
	}
	if s.Payload == "" {
		return errors.New("ciphertext is empty")
	}
	return nil
}

func (s *AESPayload) getCiphertext() (string, error) {
	if err := s.validateCiphertext(); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:v%d:%s", s.Pref, s.Version, s.Payload), nil
}

func (s *KeyModel) Update(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Updating KeyModel", "data", data)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Model(s).Update(data).Error
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
	err = db.Where(condition).Preload("Keys").First(&model).Error
	l.Debug("Finished retrieval of the KeyModel from the DB")
	return model, err
}

func FindManyKeys() ([]KeyModel, int64, error) {
	var models []KeyModel
	var count int64
	l, err := common.GetLogger()
	if err != nil {
		return models, count, err
	}
	l.Debug("Starting retrieval of the all KeyModels from the DB")
	db, err := common.GetDB()
	if err != nil {
		return models, count, err
	}
	res := db.Find(&models)
	count = res.RowsAffected
	err = res.Error
	return models, count, err
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
