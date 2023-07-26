package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
