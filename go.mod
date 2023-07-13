module github.com/miknikif/vault-auto-unseal

go 1.20

replace github.com/miknikif/vault-auto-unseal => ./command

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/render v1.0.2
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.17
)

require (
	github.com/ajg/form v1.5.1 // indirect
)
