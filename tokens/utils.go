package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"github.com/miknikif/vault-auto-unseal/policies"
)

func NewToken(tokenType string) (string, error) {
	if tokenType != TOKEN_TYPE_SERVICE && tokenType != TOKEN_TYPE_BATCH {
		return "", fmt.Errorf("token type should be one of: %s, %s", TOKEN_TYPE_SERVICE, TOKEN_TYPE_BATCH)
	}

	pref := "hvs."
	l := tokenTypeToLen[tokenType]

	if tokenType == TOKEN_TYPE_BATCH {
		pref = "hvb."
	}

	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := fmt.Sprintf("%s%s", pref, base64.StdEncoding.EncodeToString(b))

	return token, nil
}

func NewAccessor() (string, error) {
	b := make([]byte, tokenTypeToLen[TOKEN_ACCESSOR])
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	accessor := base64.StdEncoding.EncodeToString(b)

	return accessor, nil
}

func GetRemainingTTL(token TokenModel) (int, error) {
	if token.CreationTTL == 0 {
		return 0, nil
	}
	return int(token.ExpireTime.Unix() - time.Now().Unix()), nil
}

func NewRootToken() *TokenModel {
	token, _ := NewToken(TOKEN_TYPE_SERVICE)
	accessor, _ := NewAccessor()
	policyModel, _ := policies.FindOnePolicy(&policies.PolicyModel{Name: "root"})
	return &TokenModel{
		TokenID:        token,
		Accessor:       accessor,
		Policies:       []policies.PolicyModel{policyModel},
		CreationTime:   time.Now(),
		CreationTTL:    0,
		ExplicitMaxTTL: 0,
		Period:         0,
	}
}

func validateOperation(c *gin.Context) (bool, error) {
	l, _ := common.GetLogger()
	tokenID := c.Request.Header.Get(common.VAULT_TOKEN_HEADER)
	if tokenID == "" {
		return false, errors.New("token must be provided")
	}
	tokenModel, err := FindOneToken(&TokenModel{TokenID: tokenID})
	if err != nil {
		return false, errors.New("Unable to verify the token")
	}

	c.Set(common.VAULT_TOKEN, tokenID)
	c.Set(common.VAULT_TOKEN_MODEL, tokenModel)
	c.Set(common.IS_ROOT, false)

	l.Trace("Attached policies", "policies", tokenModel.Policies)
	hclPolicies := []policies.HCLPolicy{}
	for _, policy := range tokenModel.Policies {
		if policy.Name == "root" {
			l.Trace("Found Root policy attached, skipping auth", "policy", policy)
			c.Set(common.IS_ROOT, true)
			return true, nil
		}
		policyModel, err := policies.FindOnePolicy(&policy)
		if err != nil {
			return false, errors.New("Unable to retrieve policy")
		}
		text, err := common.DecFromB64(policyModel.Text)
		if err != nil {
			return false, errors.New("Unable to decode policy text")
		}
		hclPolicy, err := policies.ParseHCLPolicy(text)
		if err != nil {
			return false, errors.New("Unable to parse attached policies")
		}
		hclPolicy.Name = policy.Name
		hclPolicies = append(hclPolicies, *hclPolicy)
	}

	c.Set(common.SESSION_POLICIES, hclPolicies)

	requestPath := common.GetRequestPath(c)
	l.Trace("Found auth token validating", "path", requestPath)
	for _, hclPolicy := range hclPolicies {
		for _, path := range hclPolicy.Paths {
			if requestPath == path.Path {
				l.Trace("Found matching policy path", "path", path.Path, "policy", hclPolicy.Name)
				list, _ := strconv.ParseBool(c.Query("list"))
				requestType := c.Request.Method
				capabilitiesBitmap := path.Permissions.CapabilitiesBitmap
				capabilities := policies.GetCapabilitiesFromBitmap(capabilitiesBitmap)
				c.Set(common.PATH_CAPABILITIES, capabilities)
				if capabilities[policies.DenyCapability] {
					return false, nil
				}
				if list {
					return capabilities[policies.ListCapability], nil
				}
				switch requestType {
				case "GET":
					return capabilities[policies.ReadCapability], nil
				case "POST":
					return capabilities[policies.UpdateCapability], nil
				case "PUT":
					return capabilities[policies.UpdateCapability], nil
				case "DELETE":
					return capabilities[policies.DeleteCapability], nil
				default:
					return false, nil
				}
			}
		}
	}

	return false, nil
}
