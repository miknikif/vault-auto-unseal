package tokens

import (
	"github.com/gin-gonic/gin"
	"time"
)

type TokenSerializer struct {
	C *gin.Context
	TokenModel
}

type TokenResponse struct {
	TokenID                   string   `json:"id,omitempty"`
	Accessor                  string   `json:"accessor"`
	CreationTime              int      `json:"creation_time"`
	CreationTTL               int      `json:"creation_ttl"`
	DisplayName               string   `json:"display_name"`
	EntityID                  string   `json:"entity_id"`
	ExpireTime                string   `json:"expire_time"`
	ExplicitMaxTTL            int      `json:"explicit_max_ttl"`
	ExternalNamespacePolicies string   `json:"external_namespace_policies"`
	IdentityPolicies          []string `json:"identity_policies"`
	Meta                      string   `json:"meta"`
	NumUses                   int      `json:"num_uses"`
	Orphan                    bool     `json:"orphan"`
	Path                      string   `json:"path"`
	Policies                  []string `json:"policies"`
	Renewable                 bool     `json:"renewable"`
	TTL                       int      `json:"ttl"`
	Type                      string   `json:"type"`
}

func (s *TokenSerializer) Response() TokenResponse {
	ip := []string{}
	p := []string{}

	for _, policy := range s.IdentityPolicies {
		ip = append(ip, policy.Name)
	}

	for _, policy := range s.Policies {
		p = append(p, policy.Name)
	}

	response := TokenResponse{
		TokenID:                   s.TokenID,
		Accessor:                  s.Accessor,
		CreationTime:              s.CreationTime,
		CreationTTL:               s.CreationTTL,
		DisplayName:               s.DisplayName,
		EntityID:                  s.EntityID,
		ExpireTime:                s.ExpireTime.UTC().Format(time.RFC3339Nano),
		ExplicitMaxTTL:            s.ExplicitMaxTTL,
		ExternalNamespacePolicies: s.ExternalNamespacePolicies,
		IdentityPolicies:          ip,
		Meta:                      s.Meta,
		NumUses:                   s.NumUses,
		Orphan:                    s.Orphan,
		Path:                      s.Path,
		Policies:                  p,
		Renewable:                 s.Renewable,
		TTL:                       s.CreationTTL - (int(time.Now().Unix()) - s.CreationTime),
		Type:                      s.Type,
	}
	return response
}
