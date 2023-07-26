package policies

import (
	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

type PolicySerializer struct {
	C *gin.Context
	PolicyModel
}

type PolicyResponse struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
	Text string `json:"policy"`
}

type PoliciesSerializer struct {
	C        *gin.Context
	Policies []PolicyModel
}

type PoliciesResponse struct {
	Policies []string `json:"keys"`
}

func (s *PolicySerializer) Response() PolicyResponse {
	pt, _ := common.DecFromB64(s.Text)

	response := PolicyResponse{
		ID:   s.ID,
		Name: s.Name,
		Text: pt,
	}
	return response
}

func (s *PoliciesSerializer) Response() PoliciesResponse {
	response := PoliciesResponse{
		Policies: []string{},
	}
	for _, policy := range s.Policies {
		response.Policies = append(response.Policies, policy.Name)
	}
	return response
}
