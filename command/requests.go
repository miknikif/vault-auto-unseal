package command

import (
	"errors"
	"net/http"
)

type EncryptWithTransitKeyRequest struct {
	Plaintext string `json:"plaintext"`
}

func (er *EncryptWithTransitKeyRequest) Bind(r *http.Request) error {
	if er.Plaintext == "" {
		return errors.New("Missing plaintext.")
	}
	return nil
}

type DecryptWithTransitKeyRequest struct {
	Ciphertext string `json:"ciphertext"`
}

func (er *DecryptWithTransitKeyRequest) Bind(r *http.Request) error {
	if er.Ciphertext == "" {
		return errors.New("Missing ciphertext.")
	}
	return nil
}
