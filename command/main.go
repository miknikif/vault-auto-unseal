package command

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CmdParams struct {
	host   string
	port   int
	dbPath string
	dbName string
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrResourceAlreadyExists(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 409,
		StatusText:     "Resource already exists.",
		ErrorText:      err.Error(),
	}
}

func ErrResourceNotExists(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Requested resource not exists.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

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

type EncryptWithTransitKeyRequest struct {
	Plaintext string `json:"plaintext"`
}

func (er *EncryptWithTransitKeyRequest) Bind(r *http.Request) error {
	if er.Plaintext == "" {
		return errors.New("Missing plaintext.")
	}
	return nil
}

type EncryptWithTransitKey struct {
	Ciphertext string `json:"ciphertext"`
	KeyVersion int64  `json:"key_version"`
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

type DecryptWithTransitKeyRequest struct {
	Ciphertext string `json:"ciphertext"`
}

func (er *DecryptWithTransitKeyRequest) Bind(r *http.Request) error {
	if er.Ciphertext == "" {
		return errors.New("Missing ciphertext.")
	}
	return nil
}

type DecryptWithTransitKey struct {
	Plaintext string `json:"plaintext"`
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

func TransitKeyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key Key
		var err error

		if name := chi.URLParam(r, "name"); name != "" {
			key, err = dbGetKey(name)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "key", &key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTransitKey(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	k, err := dbGetKey(name)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, ErrResourceNotExists(errors.New(fmt.Sprintf("Transit key %s does not exist", name))))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	if err := render.Render(w, r, NewKeyResponse(&k)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func createTransitKey(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	k, err := dbGetKey(name)
	if k.Name == name {
		render.Render(w, r, ErrResourceAlreadyExists(errors.New(fmt.Sprintf("AES Key %s already exists", name))))
		return
	} else if err != nil && err != sql.ErrNoRows {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	aeskey, err := generateAESKey(32)
	if err != nil {
		render.Render(w, r, ErrRender(errors.New("unable to generate AES key")))
		return
	}

	ik := InternalKey{name: name, version: 1, aesKey: aeskey}
	err = addAESKeyToDB(ik)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	k, err = dbGetKey(ik.name)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err := render.Render(w, r, NewKeyResponse(&k)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func encryptData(key InternalKey, ek *EncryptWithTransitKeyRequest) (string, error) {
	ct, err := encryptDataWithAES(key, ek.Plaintext)

	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(ct))
	return fmt.Sprintf("vault:v1:%s", b64), nil
}

func decryptData(key InternalKey, dk *DecryptWithTransitKeyRequest) (string, error) {
	ct_data := strings.Split(dk.Ciphertext, ":")
	b64 := ct_data[len(ct_data)-1]
	ct, err := base64.StdEncoding.DecodeString(b64)

	if err != nil {
		return "", err
	}

	pt, err := decryptDataWithAES(key, string(ct))

	if err != nil {
		return "", err
	}

	return string(pt), nil
}

func encryptWithTransitKey(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value("key").(*Key)

	data := &EncryptWithTransitKeyRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	ct, err := encryptData(key.keys[0], data)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	ek := &EncryptWithTransitKey{KeyVersion: 1, Ciphertext: ct}

	render.Render(w, r, NewEncryptWithTransitKeyResponse(ek))
}

func decryptWithTransitKey(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value("key").(*Key)

	data := &DecryptWithTransitKeyRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	pt, err := decryptData(key.keys[0], data)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	dk := &DecryptWithTransitKey{Plaintext: pt}

	render.Render(w, r, NewDecryptWithTransitKeyResponse(dk))
}

func readParams(args []string) CmdParams {
	env_host := os.Getenv("VAULT_AUTO_UNSEAL_HOST")
	if env_host == "" {
		env_host = "127.0.0.1"
	}
	env_port := os.Getenv("VAULT_AUTO_UNSEAL_PORT")
	if env_port == "" {
		env_port = "8200"
	}
	env_db_path := os.Getenv("VAULT_AUTO_UNSEAL_DB_PATH")
	if env_db_path == "" {
		env_db_path = "."
	}
	env_db_name := os.Getenv("VAULT_AUTO_UNSEAL_DB_NAME")
	if env_db_name == "" {
		env_db_name = "vault-auto-unseal.db"
	}
	int_env_port, err := strconv.Atoi(env_port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return CmdParams{host: env_host, port: int_env_port, dbPath: env_db_path, dbName: env_db_name}
}

func Run(args []string) int {
	p := readParams(args)
	setDBConf(p.dbPath, p.dbName)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/transit", func(r chi.Router) {
			r.Route("/keys/{name}", func(r chi.Router) {
				r.Get("/", getTransitKey)
				r.Post("/", createTransitKey)
			})
			r.Route("/encrypt/{name}", func(r chi.Router) {
				r.Use(TransitKeyCtx)
				r.Put("/", encryptWithTransitKey)
			})
			r.Route("/decrypt/{name}", func(r chi.Router) {
				r.Use(TransitKeyCtx)
				r.Put("/", decryptWithTransitKey)
			})
		})
	})

	fmt.Printf("Listening on %s:%d\n", p.host, p.port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", p.host, p.port), r)

	return 0
}
