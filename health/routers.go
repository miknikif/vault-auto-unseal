package health

import (
	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"net/http"
)

func HealthRegister(router *gin.RouterGroup) {
	router.GET("/liveness", LivenessCheck)
	router.GET("/readiness", ReadinessCheck)
}

// LivenessCheck
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ReadinessCheck
func ReadinessCheck(c *gin.Context) {
	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("status", err))
		return
	}
	err = db.DB().Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("status", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
