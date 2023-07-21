package keys

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/miknikif/vault-auto-unseal/common"
)

type KeyModelValidator struct {
	Key struct {
		AllowPlaintextBackup bool        `json:"allow_plaintext_backup"`
		AutoRotatePeriod     int         `json:"auto_rotate_period"`
		DeletionAllowed      bool        `json:"deletion_allowed"`
		Derived              bool        `json:"derived"`
		Exportable           bool        `json:"exportable"`
		ImportedKey          bool        `json:"imported_key"`
		LatestVersion        int         `json:"latest_version"`
		MinAvailableVersion  int         `json:"min_available_version"`
		MinDecryptionVersion int         `json:"min_decryption_version"`
		MinEncryptionVersion int         `json:"min_encryption_version"`
		Name                 string      `json:"name"`
		SupportsDecryption   bool        `json:"supports_decryption"`
		SupportsDerivation   bool        `json:"supports_derivation"`
		SupportsEncryption   bool        `json:"supports_encryption"`
		SupportsSigning      bool        `json:"supports_signing"`
		Type                 KeyType     `json:"type"`
		Keys                 map[int]int `json:"keys"`
	} `json:"key"`
	keyModel KeyModel `json:"-"`
}

func (s *KeyModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}

	s.keyModel.Name = slug.Make(s.Key.Name)
	s.keyModel.AllowPlaintextBackup = s.Key.AllowPlaintextBackup
	s.keyModel.AutoRotatePeriod = s.Key.AutoRotatePeriod
	s.keyModel.DeletionAllowed = s.Key.DeletionAllowed
	s.keyModel.Derived = s.Key.Derived
	s.keyModel.Exportable = s.Key.Exportable
	s.keyModel.ImportedKey = s.Key.ImportedKey
	s.keyModel.LatestVersion = s.Key.LatestVersion
	s.keyModel.MinAvailableVersion = s.Key.MinAvailableVersion
	s.keyModel.MinDecryptionVersion = s.Key.MinDecryptionVersion
	s.keyModel.MinEncryptionVersion = s.Key.MinEncryptionVersion
	s.keyModel.Name = s.Key.Name
	s.keyModel.SupportsDecryption = s.Key.SupportsDecryption
	s.keyModel.SupportsDerivation = s.Key.SupportsDerivation
	s.keyModel.SupportsEncryption = s.Key.SupportsEncryption
	s.keyModel.SupportsSigning = s.Key.SupportsSigning
	s.keyModel.Type = s.Key.Type
	return nil
}

func NewKeyModelValidator() KeyModelValidator {
	return KeyModelValidator{}
}

func NewKeyModelValidatorFillWith(keyModel KeyModel) KeyModelValidator {
	keyModelValidator := NewKeyModelValidator()
	keyModelValidator.Key.Name = keyModel.Name
	keyModelValidator.Key.AllowPlaintextBackup = keyModel.AllowPlaintextBackup
	keyModelValidator.Key.AutoRotatePeriod = keyModel.AutoRotatePeriod
	keyModelValidator.Key.DeletionAllowed = keyModel.DeletionAllowed
	keyModelValidator.Key.Derived = keyModel.Derived
	keyModelValidator.Key.Exportable = keyModel.Exportable
	keyModelValidator.Key.ImportedKey = keyModel.ImportedKey
	keyModelValidator.Key.LatestVersion = keyModel.LatestVersion
	keyModelValidator.Key.MinAvailableVersion = keyModel.MinAvailableVersion
	keyModelValidator.Key.MinDecryptionVersion = keyModel.MinDecryptionVersion
	keyModelValidator.Key.MinEncryptionVersion = keyModel.MinEncryptionVersion
	keyModelValidator.Key.SupportsDecryption = keyModel.SupportsDecryption
	keyModelValidator.Key.SupportsDerivation = keyModel.SupportsDerivation
	keyModelValidator.Key.SupportsEncryption = keyModel.SupportsEncryption
	keyModelValidator.Key.SupportsSigning = keyModel.SupportsSigning
	keyModelValidator.Key.Type = keyModel.Type
	return keyModelValidator
}
