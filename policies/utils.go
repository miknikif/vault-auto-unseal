package policies

import "github.com/miknikif/vault-auto-unseal/common"

func NewRootPolicy() *PolicyModel {
	policyText := `path "*" {
    capabilities = ["read", "create", "list", "update", "delete", "sudo"]
}`
	return &PolicyModel{
		Name: "root",
		Text: common.EncToB64(policyText),
	}
}

func NewDefaultPolicy() *PolicyModel {
	policyText := `path "auth/token/self-lookup" {
    capabilities = ["read"]
}`
	return &PolicyModel{
		Name: "default",
		Text: common.EncToB64(policyText),
	}
}
