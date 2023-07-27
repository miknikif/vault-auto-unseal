package tokens

import (
	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

const (
	VAULT_TOKEN = "vaultToken"
)

func AuthMiddleware() gin.HandlerFunc {
	l, _ := common.GetLogger()
	return func(c *gin.Context) {
		l.Debug("Running AuthMiddleware")
		l.Trace("Read X-Vault-Auth header")
		token := c.Request.Header.Get("X-Vault-Token")
		l.Trace("Found auth token", "token", token)
		c.Set(VAULT_TOKEN, token)
		c.Next()
	}
}
