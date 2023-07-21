package keys

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type KeySerializer struct {
	C *gin.Context
	KeyModel
}

type KeyResponse struct {
	AllowPlaintextBackup bool    `json:"allow_plaintext_backup"`
	AutoRotatePeriod     int     `json:"auto_rotate_period"`
	DeletionAllowed      bool    `json:"deletion_allowed"`
	Derived              bool    `json:"derived"`
	Exportable           bool    `json:"exportable"`
	ImportedKey          bool    `json:"imported_key"`
	LatestVersion        int     `json:"latest_version"`
	MinAvailableVersion  int     `json:"min_available_version"`
	MinDecryptionVersion int     `json:"min_decryption_version"`
	MinEncryptionVersion int     `json:"min_encryption_version"`
	Name                 string  `json:"name"`
	SupportsDecryption   bool    `json:"supports_decryption"`
	SupportsDerivation   bool    `json:"supports_derivation"`
	SupportsEncryption   bool    `json:"supports_encryption"`
	SupportsSigning      bool    `json:"supports_signing"`
	Type                 KeyType `json:"type"`
	// Keys                 map[int]int `json:"keys"`
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
		Name:                 slug.Make(s.Name),
		SupportsDecryption:   s.SupportsDecryption,
		SupportsDerivation:   s.SupportsDerivation,
		SupportsEncryption:   s.SupportsEncryption,
		SupportsSigning:      s.SupportsSigning,
		Type:                 s.Type,
		//Keys:                 s.Keys,
	}
	return response
}

// func (s *KeySerializer) Response() []KeyResponse {
//     response := []KeyResponse{}
//     for _, key := range s.Keys {
//         serializer := KeySerializer{s.C, key}
//     }
// }
