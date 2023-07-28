package keys

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

type KeyModelValidator struct {
	AllowPlaintextBackup bool     `json:"allow_plaintext_backup"`
	AutoRotatePeriod     string   `json:"auto_rotate_period"`
	DeletionAllowed      string   `json:"deletion_allowed"`
	Derived              string   `json:"derived"`
	Exportable           string   `json:"exportable"`
	MinDecryptionVersion string   `json:"min_decryption_version"`
	MinEncryptionVersion string   `json:"min_encryption_version"`
	Name                 string   `json:"-"`
	Type                 KeyType  `json:"type"`
	keyModel             KeyModel `json:"-"`
}

type EncryptDataValidator struct {
	Plaintext  string `json:"plaintext"`
	aesPayload AESPayload
}

type DecryptDataValidator struct {
	Ciphertext string `json:"ciphertext"`
	aesPayload AESPayload
}

func (s *KeyModelValidator) Bind(c *gin.Context) error {
	l, _ := common.GetLogger()
	err := common.Bind(c, s)
	if err != nil {
		return err
	}

	s.keyModel.AllowPlaintextBackup = false
	s.keyModel.AutoRotatePeriod = 0
	s.keyModel.DeletionAllowed = false
	s.keyModel.Derived = false
	s.keyModel.Exportable = false
	s.keyModel.ImportedKey = false
	s.keyModel.MinAvailableVersion = 1
	s.keyModel.MinDecryptionVersion = 1
	s.keyModel.MinEncryptionVersion = 1
	s.keyModel.SupportsDecryption = true
	s.keyModel.SupportsDerivation = false
	s.keyModel.SupportsEncryption = true
	s.keyModel.SupportsSigning = false
	s.keyModel.Type = KEY_TYPE_AES256_GCM96

	autoRotatePeriod := common.ParseInt(s.AutoRotatePeriod, 0)
	minDecryptionVersion := common.ParseInt(s.MinDecryptionVersion, 1)
	minEncryptionVersion := common.ParseInt(s.MinEncryptionVersion, 1)
	deletionAllowed := common.ParseBool(s.DeletionAllowed, false)
	derived := common.ParseBool(s.Derived, false)
	exportable := common.ParseBool(s.Exportable, false)

	s.keyModel.AllowPlaintextBackup = s.AllowPlaintextBackup
	s.keyModel.AutoRotatePeriod = autoRotatePeriod
	s.keyModel.DeletionAllowed = deletionAllowed
	s.keyModel.Derived = derived
	s.keyModel.Exportable = exportable
	s.keyModel.MinDecryptionVersion = minDecryptionVersion
	s.keyModel.MinEncryptionVersion = minEncryptionVersion
	s.keyModel.Name = s.Name

	if s.keyModel.Keys == nil || len(s.keyModel.Keys) < 1 {
		key, err := createNewAESKeyModel(1, AES_KEY_SIZE_256)
		if err != nil {
			return err
		}

		s.keyModel.Keys = []AESKeyModel{key}
	}

	s.keyModel.LatestVersion = len(s.keyModel.Keys)
	if minDecryptionVersion > len(s.keyModel.Keys) || minEncryptionVersion > len(s.keyModel.Keys) {
		return errors.New("MinEncryptionVersion and MinDecryptionVersion are referencing not existing keys")
	}

	l.Debug("KeyModelValidator", "keys_amount", len(s.keyModel.Keys), "keyModel", s.keyModel)

	return nil
}

func NewKeyModelValidator() KeyModelValidator {
	return KeyModelValidator{}
}

func NewKeyModelValidatorFillWith(keyModel KeyModel) KeyModelValidator {
	keyModelValidator := NewKeyModelValidator()
	keyModelValidator.Name = keyModel.Name
	keyModelValidator.AllowPlaintextBackup = keyModel.AllowPlaintextBackup
	keyModelValidator.AutoRotatePeriod = fmt.Sprint(keyModel.AutoRotatePeriod)
	keyModelValidator.DeletionAllowed = fmt.Sprint(keyModel.DeletionAllowed)
	keyModelValidator.Derived = fmt.Sprint(keyModel.Derived)
	keyModelValidator.Exportable = fmt.Sprint(keyModel.Exportable)
	keyModelValidator.MinDecryptionVersion = fmt.Sprint(keyModel.MinDecryptionVersion)
	keyModelValidator.MinEncryptionVersion = fmt.Sprint(keyModel.MinEncryptionVersion)
	keyModelValidator.Type = keyModel.Type
	keyModelValidator.keyModel.Keys = keyModel.Keys
	return keyModelValidator
}

func NewEncryptDataValidator() EncryptDataValidator {
	return EncryptDataValidator{}
}

func (s *EncryptDataValidator) Bind(c *gin.Context) error {
	l, _ := common.GetLogger()
	l.Debug("EncryptDataValidator.Bind - start")
	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	if err := s.Validate(); err != nil {
		return err
	}
	l.Debug("EncryptDataValidator.Bind - end")
	return nil
}

func (s *EncryptDataValidator) Validate() error {
	l, _ := common.GetLogger()
	l.Debug("EncryptDataValidator.Validate - start")
	s.aesPayload.Plaintext = s.Plaintext

	if err := s.aesPayload.validatePlaintext(); err != nil {
		return err
	}

	l.Debug("EncryptDataValidator.Validate - end")
	return nil
}

func NewEncryptDataValidatorFillWith(aesPayload AESPayload) EncryptDataValidator {
	encryptDataValidator := NewEncryptDataValidator()
	encryptDataValidator.Plaintext = aesPayload.Plaintext
	return encryptDataValidator
}

func NewDecryptDataValidator() DecryptDataValidator {
	return DecryptDataValidator{}
}

func (s *DecryptDataValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}

	if s.Ciphertext == "" {
		return errors.New("ciphertext should be specified")
	}

	data := strings.Split(s.Ciphertext, ":")

	if len(data) != 3 {
		return errors.New("Wrong format of the ciphertext is recieved")
	}

	s.aesPayload.Pref = data[0]
	s.aesPayload.Version = common.ParseInt(string(data[1][1]), 1)
	s.aesPayload.Payload = data[2]

	if err := s.aesPayload.validateCiphertext(); err != nil {
		return err
	}

	return nil
}
