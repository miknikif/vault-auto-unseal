package keys

import (
	"github.com/gin-gonic/gin"
)

type KeySerializer struct {
	C *gin.Context
	KeyModel
}

type EncryptDataSerializer struct {
	C *gin.Context
	AESPayload
}

type DecryptDataSerializer struct {
	C *gin.Context
	AESPayload
}

type EncryptDataResponse struct {
	Ciphertext string `json:"ciphertext"`
	Version    int    `json:"version"`
}

type DecryptDataResponse struct {
	Plaintext string `json:"plaintext"`
}

type KeyResponse struct {
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
}

type KeysSerializer struct {
	C    *gin.Context
	Keys []KeyModel
}

type KeysResponse struct {
	Keys []string `json:"keys"`
}

func (s *KeysSerializer) Response() KeysResponse {
	response := KeysResponse{
		Keys: []string{},
	}
	for _, key := range s.Keys {
		response.Keys = append(response.Keys, key.Name)
	}
	return response
}

func (s *KeySerializer) Response() KeyResponse {
	response := KeyResponse{
		AllowPlaintextBackup: s.AllowPlaintextBackup,
		AutoRotatePeriod:     s.AutoRotatePeriod,
		DeletionAllowed:      s.DeletionAllowed,
		Derived:              s.Derived,
		Exportable:           s.Exportable,
		ImportedKey:          s.ImportedKey,
		LatestVersion:        s.LatestVersion,
		MinAvailableVersion:  s.MinAvailableVersion,
		MinDecryptionVersion: s.MinDecryptionVersion,
		MinEncryptionVersion: s.MinEncryptionVersion,
		Name:                 s.Name,
		SupportsDecryption:   s.SupportsDecryption,
		SupportsDerivation:   s.SupportsDerivation,
		SupportsEncryption:   s.SupportsEncryption,
		SupportsSigning:      s.SupportsSigning,
		Type:                 s.Type,
		Keys:                 make(map[int]int),
	}

	for _, key := range s.Keys {
		response.Keys[key.Version] = key.Name
	}

	return response
}

func (s *EncryptDataSerializer) Response() EncryptDataResponse {
	ct, _ := s.AESPayload.getCiphertext()
	response := EncryptDataResponse{
		Ciphertext: ct,
		Version:    s.AESPayload.Version,
	}
	return response
}

func (s *DecryptDataSerializer) Response() DecryptDataResponse {
	response := DecryptDataResponse{
		Plaintext: s.AESPayload.Plaintext,
	}
	return response
}
