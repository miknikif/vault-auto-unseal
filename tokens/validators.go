package tokens

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"github.com/miknikif/vault-auto-unseal/policies"
)

type TokenLookupModelValidator struct {
	TokenID          string     `form:"token" json:"token"`
	Accessor         string     `form:"accessor" json:"accessor"`
	findWithAccessor bool       `json:"-"`
	isSelf           bool       `json:"-"`
	tokenModel       TokenModel `json:"-"`
}

func (s *TokenLookupModelValidator) Bind(c *gin.Context) error {
	if !s.isSelf {
		if err := common.Bind(c, s); err != nil {
			return err
		}
	} else {
		s.TokenID = c.GetString(VAULT_TOKEN)
	}

	s.tokenModel.TokenID = s.TokenID
	s.tokenModel.Accessor = s.Accessor

	if s.tokenModel.TokenID == "" && s.tokenModel.Accessor == "" {
		return errors.New("token ID or accessor must be specified")
	}

	s.findWithAccessor = s.tokenModel.TokenID == ""

	return nil
}

func NewTokenLookupModelValidator(self bool) TokenLookupModelValidator {
	tokenLookupModelValidator := TokenLookupModelValidator{
		isSelf: self,
	}
	return tokenLookupModelValidator
}

func NewTokenLookupModelValidatorFillWith(tokenModel TokenModel) TokenLookupModelValidator {
	tokenLookupModelValidator := NewTokenLookupModelValidator(false)
	tokenLookupModelValidator.TokenID = tokenModel.TokenID
	tokenLookupModelValidator.Accessor = tokenModel.Accessor
	return tokenLookupModelValidator
}

type TokenModelValidator struct {
	Policies       []string   `json:"policies"`
	TTL            string     `json:"ttl"`
	ExplicitMaxTTL string     `json:"explicit_max_ttl"`
	Period         string     `json:"period"`
	DisplayName    string     `json:"display_name"`
	NumUses        int        `json:"num_uses"`
	Renewable      bool       `json:"renewable"`
	Type           string     `json:"type"`
	EntityAlias    string     `json:"entity_alias"`
	tokenModel     TokenModel `json:"-"`
}

func (s *TokenModelValidator) Bind(c *gin.Context) error {
	if err := common.Bind(c, s); err != nil {
		return err
	}

	p := []policies.PolicyModel{}
	for _, policy := range s.Policies {
		pol, err := policies.FindOnePolicy(&policies.PolicyModel{Name: policy})
		if err != nil {
			return fmt.Errorf("Policy %s not found", policy)
		}
		p = append(p, pol)
	}

	ttl, err := time.ParseDuration(s.TTL)
	if err != nil {
		return errors.New("unable to parse ttl")
	}
	explicitMaxTTL, err := time.ParseDuration(s.ExplicitMaxTTL)
	if err != nil {
		return errors.New("unable to parse explicit max ttl")
	}
	period, err := time.ParseDuration(s.Period)
	if err != nil {
		return errors.New("unable to parse period")
	}
	if strings.ToLower(s.Type) != TOKEN_TYPE_BATCH && strings.ToLower(s.Type) != TOKEN_TYPE_SERVICE {
		return fmt.Errorf("token type should be one of the following: %s or %s", TOKEN_TYPE_SERVICE, TOKEN_TYPE_BATCH)
	}
	if int(ttl) == 0 && int(period) == 0 {
		return fmt.Errorf("ttl and period cannot be 0")
	} else if int(ttl) == 0 && int(period) > 0 {
		ttl = period
	}
	if int(explicitMaxTTL) != 0 && int(period) > 0 {
		return fmt.Errorf("explicit_max_ttl can be only 0 if period is specified")
	}

	s.tokenModel.Policies = p
	s.tokenModel.CreationTTL = int(ttl.Seconds())
	s.tokenModel.ExplicitMaxTTL = int(explicitMaxTTL.Seconds())
	s.tokenModel.Period = int(period.Seconds())
	s.tokenModel.DisplayName = s.DisplayName
	s.tokenModel.NumUses = s.NumUses
	s.tokenModel.Renewable = s.Renewable
	s.tokenModel.Type = strings.ToLower(s.Type)
	s.tokenModel.EntityID = s.EntityAlias
	s.tokenModel.CreationTime = time.Now()
	s.tokenModel.Path = c.FullPath()
	s.tokenModel.ExpireTime = time.Now().Add(ttl)

	if s.tokenModel.TokenID == "" {
		token, err := NewToken(s.tokenModel.Type)
		if err != nil {
			return err
		}
		s.tokenModel.TokenID = token
	}

	if s.tokenModel.Accessor == "" {
		accessor, err := NewAccessor()
		if err != nil {
			return err
		}
		s.tokenModel.Accessor = accessor
	}

	return nil
}

func NewTokenModelValidator() TokenModelValidator {
	tokenModelValidator := TokenModelValidator{}
	return tokenModelValidator
}

func NewTokenModelValidatorFillWith(tokenModel TokenModel) TokenModelValidator {
	tokenModelValidator := NewTokenModelValidator()

	p := []string{}
	for _, policy := range tokenModel.Policies {
		p = append(p, policy.Name)
	}

	tokenModelValidator.Policies = p
	tokenModelValidator.TTL = fmt.Sprintf("%ds", tokenModel.TTL)
	tokenModelValidator.ExplicitMaxTTL = fmt.Sprintf("%ds", tokenModel.ExplicitMaxTTL)
	tokenModelValidator.Period = fmt.Sprintf("%ds", tokenModel.Period)
	tokenModelValidator.DisplayName = tokenModel.DisplayName
	tokenModelValidator.NumUses = tokenModel.NumUses
	tokenModelValidator.Renewable = tokenModel.Renewable
	tokenModelValidator.Type = tokenModel.Type
	tokenModelValidator.EntityAlias = tokenModel.EntityID
	return tokenModelValidator
}
