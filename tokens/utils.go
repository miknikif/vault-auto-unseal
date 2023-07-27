package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

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
	return &TokenModel{
		TokenID:  token,
		Accessor: accessor,
		Policies: []policies.PolicyModel{
			{
				Name: "root",
			},
		},
		CreationTime:   time.Now(),
		CreationTTL:    0,
		ExplicitMaxTTL: 0,
		Period:         0,
	}
}
