package tokens

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
)

func AuthMiddleware() gin.HandlerFunc {
	l, _ := common.GetLogger()
	return func(c *gin.Context) {
		l.Debug("Running AuthMiddleware")
		l.Trace("Read X-Vault-Token header")

		res, err := validateOperation(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, common.NewError("auth", err))
			return
		}
		if !res {
			c.AbortWithStatusJSON(http.StatusForbidden, common.NewError("auth", errors.New("permission denied")))
			return
		}

		c.Next()
	}
}
