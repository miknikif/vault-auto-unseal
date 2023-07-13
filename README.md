# Vault transit auto unseal

**Warning**: Before using anything from this repo, consider the following:
- Implementation isn't checked from the security perspective
- A single key type supported - AES256
- This app is handling only some of the transit secret engine APIs
- Code quality isn't good (right now at least)
- Encryption AES256 key is stored unencrypted in the sqlite db (alongside with the app)
- Whole app was created to train with GO
- Probably the primary usage of this - homelab, or testing dev vault environment

## What this

Vault has an awesome feature of auto unsealing it using the transit secrets engine. The problem is - you need a whole separate vault cluster with really high uptime and realability.
This small go app is mocking transit secret engine and implements endpoints required to make auto unseal work.

## How to use

1. Container
	- Clone this repo
	- Build the container image using something like: `docker build -rm -t <image_name> -f Dockerfile .`
	- By default server will use port `8200` and db will be saved at the following path: `/w/db/vaseal.db`
	- Run builded image: `docker run -itp 8200:8200 -v <image_name>`
	- Create a new key by running `curl -X POST http://localhost:8200/v1/transit/keys/<key_name>`
	- Apply the following example config to your vault:
	    ```hcl
		seal "transit" {
          address            = "http://<ip>:8200"
          token              = "s.Qf1s5zigZ4OX6akYjQXJC1jY" # any random token will work, it's not used
          disable_renewal    = "true" # disable vault from trying to lookup and renew the token
          key_name           = "<key_name>" 
          mount_path         = "transit/"
          tls_skip_verify    = "true" # I've not implemented TLS
        }
		```
	- Start vault and init it - `vault operator init`
2. Binary
	- Download precompiled binary, or compile it yourself:
		- To compile it yourself clone the repo
		- Install go
		- Run `go mod download`
		- Compile binary `CGO_ENABLED=1 go build -ldflags="-w -s" -o /vault-auto-unseal main.go`
	- By default server will use port `8200` and db will be saved at the following path: `./vault-auto-unseal.db`
	- Run compiled binary: `./vault-auto-unseal`
	- Create a new key by running `curl -X POST http://localhost:8200/v1/transit/keys/<key_name>`
	- Apply the following example config to your vault:
	    ```hcl
		seal "transit" {
          address            = "http://<ip>:8200"
          token              = "s.Qf1s5zigZ4OX6akYjQXJC1jY" # any random token will work, it's not used
          disable_renewal    = "true" # disable vault from trying to lookup and renew the token
          key_name           = "<key_name>" 
          mount_path         = "transit/"
          tls_skip_verify    = "true" # I've not implemented TLS
        }
		```
	- Start vault and init it - `vault operator init`

### ENV Vars
1. VAULT_AUTO_UNSEAL_HOST - `string` (default: `0.0.0.0`)
1. VAULT_AUTO_UNSEAL_PORT - `int` (default: `8200`)
1. VAULT_AUTO_UNSEAL_DB_PATH - `string` (default: `.`)
1. VAULT_AUTO_UNSEAL_DB_NAME - `string` (default: `vault-auto-unseal.db`)

`VAULT_AUTO_UNSEAL_DB_PATH` and `VAULT_AUTO_UNSEAL_DB_NAME` are building the os path, so by default it'll create a DB on the following path `./vault-auto-unseal.db`
