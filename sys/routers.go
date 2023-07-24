package sys

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/miknikif/vault-auto-unseal/common"
	"net/http"
)

func HealthRegister(router *gin.RouterGroup) {
	router.GET("/liveness", LivenessCheck)
	router.GET("/readiness", ReadinessCheck)
	router.GET("/health", HealthRetrieve)
	router.GET("/seal-status", SealStatusRetrieve)
	router.GET("/leader", LeaderStatusRetrieve)
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

// Using same status format as original vault is using
func HealthRetrieve(c *gin.Context) {
	healthModel, err := GetHealthStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("health", errors.New("Unable to get health")))
		return
	}
	serializer := HealthSerializer{c, healthModel}
	c.JSON(http.StatusOK, serializer.Response())
}

// Using same status format as original vault is using
func SealStatusRetrieve(c *gin.Context) {
	healthModel, err := GetSealStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("health", errors.New("Unable to get seal-status")))
		return
	}
	serializer := SealStatusSerializer{c, healthModel}
	c.JSON(http.StatusOK, serializer.Response())
}

// Using same status format as original vault is using
func LeaderStatusRetrieve(c *gin.Context) {
	leaderStatusModel, err := GetLeaderStatus()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("health", errors.New("Unable to get leader status")))
		return
	}
	serializer := LeaderStatusSerializer{c, leaderStatusModel}
	c.JSON(http.StatusOK, serializer.Response())
}
