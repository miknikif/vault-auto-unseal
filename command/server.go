package command

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"github.com/miknikif/vault-auto-unseal/health"
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

	if c.Args.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	v1 := router.Group("/v1")
	health.HealthRegister(v1.Group("/"))
	keys.KeysRegister(v1.Group("/transit/keys"))

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", c.Args.Host, c.Args.Port),
		Handler:  router,
		ErrorLog: c.Logger.StandardLogger(nil),
	}

	if c.TLS != nil {
		server.TLSConfig = c.TLS.TLSConfig
		fmt.Println(c.TLS.BundleCrt)
		c.Logger.Info(fmt.Sprintf("Starting HTTPS server at https://%s:%d", c.Args.Host, c.Args.Port))
		err = server.ListenAndServeTLS(c.TLS.BundleCrt, c.TLS.TLSKey)
		if err != nil {
			return err
		}
	} else {
		c.Logger.Info(fmt.Sprintf("Starting HTTP server at http://%s:%d", c.Args.Host, c.Args.Port))
		err = server.ListenAndServe()
		if err != nil {
			return err
		}
	}
	return nil
}
