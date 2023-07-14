package command

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SuccessResponse struct {
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
}

func (res *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.HTTPStatusCode)
	return nil
}

var ResponseResourceDeleted = &SuccessResponse{HTTPStatusCode: 200, StatusText: "OK"}

type GenericResponseData interface{}

type GenericResponse struct {
	RequestID     string              `json:"request_id"`
	LeaseID       string              `json:"lease_id"`
	LeaseDuration int                 `json:"lease_duration"`
	Renewable     bool                `json:"renewable"`
	Data          GenericResponseData `json:"data"`
	Warnings      *[]string           `json:"warnings"`
}

func (gr *GenericResponse) Render(w http.ResponseWriter, r *http.Request) error {
	gr.RequestID = uuid.New().String()
	gr.LeaseID = ""
	gr.LeaseDuration = 0
	gr.Renewable = false
	return nil
}

type KeyResponse struct {
	*GenericResponse
	Data *Key `json:"data"`
}

func (kr *KeyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	keys := map[int]int{}
	for i, v := range kr.Data.keys {
		keys[i+1] = v.id
	}
	kr.Data.Keys = keys
	return nil
}

func NewKeyResponse(key *Key) *KeyResponse {
	resp := &KeyResponse{Data: key, GenericResponse: &GenericResponse{}}
	return resp
}

type EncryptWithTransitKey struct {
	Ciphertext string `json:"ciphertext"`
	KeyVersion int64  `json:"key_version"`
}

type DecryptWithTransitKey struct {
	Plaintext string `json:"plaintext"`
}

type EncryptWithTransitKeyResponse struct {
	*GenericResponse
	Data *EncryptWithTransitKey `json:"data"`
}

func (er *EncryptWithTransitKeyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	er.Data.KeyVersion = 1
	return nil
}

func NewEncryptWithTransitKeyResponse(ek *EncryptWithTransitKey) *EncryptWithTransitKeyResponse {
	resp := &EncryptWithTransitKeyResponse{Data: ek, GenericResponse: &GenericResponse{}}
	return resp
}

type DecryptWithTransitKeyResponse struct {
	*GenericResponse
	Data *DecryptWithTransitKey `json:"data"`
}

func (er *DecryptWithTransitKeyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewDecryptWithTransitKeyResponse(dk *DecryptWithTransitKey) *DecryptWithTransitKeyResponse {
	resp := &DecryptWithTransitKeyResponse{Data: dk, GenericResponse: &GenericResponse{}}
	return resp
}
