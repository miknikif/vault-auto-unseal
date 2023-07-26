package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	l, _ := GetLogger()
	return func(c *gin.Context) {
		l.Debug("Running RequestIDMiddleware")
		requestID, _ := uuid.NewRandom()
		c.Set("request_id", requestID.String())
		l.Debug(fmt.Sprintf("Set RequestIDMiddleware:request_id to %s", requestID.String()))
	}
}

func JSONMiddleware(replaceExistingContentType bool) gin.HandlerFunc {
	l, _ := GetLogger()
	return func(c *gin.Context) {
		l.Debug("Running JSONMiddleware")
		if replaceExistingContentType {
			l.Trace("Override Content-Type to application/json")
			c.Request.Header.Set("Content-Type", "application/json")
			c.Next()
			return
		} else {
			ct := c.ContentType()
			if ct == "" {
				l.Trace("Set Content-Type to application/json")
				c.Request.Header.Set("Content-Type", "application/json")
				c.Next()
				return
			}
		}
		l.Trace("Content-Type header is left untact")
		c.Next()
	}
}
