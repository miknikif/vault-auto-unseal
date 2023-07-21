package command

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"github.com/miknikif/vault-auto-unseal/keys"
)

// Migrate provided DB
func Migrate(c *common.Config) {
	c.Logger.Info(fmt.Sprintf("Migrating %s", c.Args.DBName))
	c.DB.AutoMigrate(&keys.AESKeyModel{})
	c.DB.AutoMigrate(&keys.KeyModel{})
}

// Start HTTP Server
func StartHttpServer() error {
	c, err := common.GetConfig()
	if err != nil {
		return err
	}
	Migrate(c)
	defer c.DB.Close()

	r := gin.Default()
	v1 := r.Group("/v1")
	keys.KeysRegister(v1.Group("/transit/keys"))
	v1.Group("/liveness").GET("", common.LivenessCheck)
	v1.Group("/readiness").GET("", common.ReadinessCheck)

	c.Logger.Info(fmt.Sprintf("Starting HTTP server at %s://%s:%d", c.Args.TLS.Proto, c.Args.Host, c.Args.Port))
	r.Run(fmt.Sprintf("%s:%d", c.Args.Host, c.Args.Port))

	return nil
}
