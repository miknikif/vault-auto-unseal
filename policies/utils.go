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

func GetCapabilitiesFromBitmap(bitmap uint32) map[string]bool {
	return map[string]bool{
		DenyCapability:   bitmap&DenyCapabilityInt > 0,
		SudoCapability:   bitmap&SudoCapabilityInt > 0,
		CreateCapability: bitmap&CreateCapabilityInt > 0,
		UpdateCapability: bitmap&UpdateCapabilityInt > 0,
		ReadCapability:   bitmap&ReadCapabilityInt > 0,
		ListCapability:   bitmap&ListCapabilityInt > 0,
		DeleteCapability: bitmap&DeleteCapabilityInt > 0,
	}
}
