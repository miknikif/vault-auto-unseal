package policies

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

type PolicyModelValidator struct {
	Name        string      `json:"-"`
	Text        string      `form:"policy" json:"policy"`
	policyModel PolicyModel `json:"-"`
}

func (s *PolicyModelValidator) Bind(c *gin.Context) error {
	if err := common.Bind(c, s); err != nil {
		return err
	}

	s.policyModel.Name = s.Name
	s.policyModel.Text = s.Text

	if s.policyModel.Name == "" || s.policyModel.Text == "" {
		return errors.New("policy - name or policy text is empty")
	}

	if _, err := parseHCLPolicy(s.Text); err != nil {
		return err
	}

	s.policyModel.Text = common.EncToB64(s.policyModel.Text)

	return nil
}

func NewPolicyModelValidator() PolicyModelValidator {
	policyModelValidator := PolicyModelValidator{}
	return policyModelValidator
}

func NewPolicyModelValidatorFillWith(policyModel PolicyModel) PolicyModelValidator {
	policyModelValidator := NewPolicyModelValidator()
	policyModelValidator.Text = policyModel.Text
	policyModelValidator.Name = policyModel.Name

	return policyModelValidator
}
