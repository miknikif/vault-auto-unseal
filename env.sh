#!/usr/bin/env bash

export HOSTNAME="localhost"
export VAULT_AUTO_UNSEAL_HOST="127.0.0.1"
export VAULT_AUTO_UNSEAL_PRODUCTION="false"

# TLS
# export VAULT_AUTO_UNSEAL_CLIENT_CA_CRT_PATH="${PWD}/tls/ca.crt"
export VAULT_AUTO_UNSEAL_CA_CRT_PATH="${PWD}/tls/ca.crt"
export VAULT_AUTO_UNSEAL_TLS_CRT_PATH="${PWD}/tls/tls.crt"
export VAULT_AUTO_UNSEAL_TLS_KEY_PATH="${PWD}/tls/tls.key"

