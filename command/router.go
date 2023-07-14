package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

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

func deleteTransitKey(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	ik, err := dbGetKey(name)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = delAESKeyFromDB(ik.keys[1])
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err := render.Render(w, r, ResponseResourceDeleted); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func encryptWithTransitKey(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value("key").(*Key)

	data := &EncryptWithTransitKeyRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	ct, err := encryptData(key.keys[1], data)

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
	pt, err := decryptData(key.keys[1], data)

	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	dk := &DecryptWithTransitKey{Plaintext: pt}

	render.Render(w, r, NewDecryptWithTransitKeyResponse(dk))
}

func createRouter() *chi.Mux {
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
				r.Delete("/", deleteTransitKey)
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

	return r
}
