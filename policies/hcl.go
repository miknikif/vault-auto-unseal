package policies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/helper/hclutil"
	"github.com/miknikif/vault-auto-unseal/common"
)

const (
	DenyCapability   = "deny"
	CreateCapability = "create"
	ReadCapability   = "read"
	UpdateCapability = "update"
	DeleteCapability = "delete"
	ListCapability   = "list"
	SudoCapability   = "sudo"
	RootCapability   = "root"
	PatchCapability  = "patch"
)

const (
	DenyCapabilityInt uint32 = 1 << iota
	CreateCapabilityInt
	ReadCapabilityInt
	UpdateCapabilityInt
	DeleteCapabilityInt
	ListCapabilityInt
	SudoCapabilityInt
	PatchCapabilityInt
)

var cap2Int = map[string]uint32{
	DenyCapability:   DenyCapabilityInt,
	CreateCapability: CreateCapabilityInt,
	ReadCapability:   ReadCapabilityInt,
	UpdateCapability: UpdateCapabilityInt,
	DeleteCapability: DeleteCapabilityInt,
	ListCapability:   ListCapabilityInt,
	SudoCapability:   SudoCapabilityInt,
	PatchCapability:  PatchCapabilityInt,
}

type ACLPermissions struct {
	CapabilitiesBitmap uint32
	AllowedParameters  map[string][]interface{}
	DeniedParameters   map[string][]interface{}
	RequiredParameters []string
}

type HCLPolicyPathRules struct {
	Path                string
	Policy              string
	Permissions         *ACLPermissions
	IsPrefix            bool
	HasSegmentWildcards bool
	Capabilities        []string
}

type HCLPolicy struct {
	Name  string                `hcl:"name"`
	Paths []*HCLPolicyPathRules `hcl:"-"`
	Raw   string
}

func parsePaths(result *HCLPolicy, list *ast.ObjectList) error {
	paths := make([]*HCLPolicyPathRules, 0, len(list.Items))
	for _, item := range list.Items {
		key := "path"
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		valid := []string{
			"comment",
			"policy",
			"capabilities",
			"allowed_parameters",
			"denied_parameters",
			"required_parameters",
		}
		if err := hclutil.CheckHCLKeys(item.Val, valid); err != nil {
			return fmt.Errorf("path %q, error: %w", key, err)
		}

		var pc HCLPolicyPathRules

		pc.Permissions = new(ACLPermissions)
		pc.Path = key

		if err := hcl.DecodeObject(&pc, item.Val); err != nil {
			return fmt.Errorf("path %q, error: %w", key, err)
		}

		// Strip a leading '/' as paths in Vault start after the / in the API path
		if len(pc.Path) > 0 && pc.Path[0] == '/' {
			pc.Path = pc.Path[1:]
		}

		if strings.Contains(pc.Path, "+*") {
			return fmt.Errorf("path %q: invalid use of wildcards ('+*' is forbidden)", pc.Path)
		}

		if pc.Path == "+" || strings.Count(pc.Path, "/+") > 0 || strings.HasPrefix(pc.Path, "+/") {
			pc.HasSegmentWildcards = true
		}

		if strings.HasSuffix(pc.Path, "*") {
			// If there are segment wildcards, don't actually strip the
			// trailing asterisk, but don't want to hit the default case
			if !pc.HasSegmentWildcards {
				// Strip the glob character if found
				pc.Path = strings.TrimSuffix(pc.Path, "*")
				pc.IsPrefix = true
			}
		}

		pc.Permissions.CapabilitiesBitmap = 0
		for _, cap := range pc.Capabilities {
			switch cap {
			case DenyCapability:
				pc.Capabilities = []string{DenyCapability}
				pc.Permissions.CapabilitiesBitmap = DenyCapabilityInt
				goto PathFinished
			case CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability, SudoCapability, PatchCapability:
				pc.Permissions.CapabilitiesBitmap |= cap2Int[cap]
			default:
				return fmt.Errorf("path %q: invalid capability %q", key, cap)
			}
		}
	PathFinished:
		paths = append(paths, &pc)
	}

	result.Paths = paths
	return nil
}

func ParseHCLPolicy(src string) (*HCLPolicy, error) {
	l, err := common.GetLogger()
	if err != nil {
		return nil, err
	}
	l.Debug("Starting policy parsing")
	root, err := hcl.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("failed to parse policy: does not contain a root object")
	}

	valid := []string{
		"name",
		"path",
	}

	if err := hclutil.CheckHCLKeys(list, valid); err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	p := HCLPolicy{
		Raw: src,
	}
	if err := hcl.DecodeObject(&p, list); err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	if o := list.Filter("path"); len(o.Items) > 0 {
		if err := parsePaths(&p, o); err != nil {
			return nil, fmt.Errorf("failed to parse policy: %w", err)
		}
	}

	l.Debug("Parsed policy", "hcl_policy", p)
	return &p, nil
}
