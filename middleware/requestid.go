package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func RequestID(c *gin.Context) {
	if requestID := c.Request.Header.Get("X-Request-Id"); requestID == "" {
		requestID := GenRequestID()
		c.Request.Header.Set("X-Request-Id", requestID)
		c.Set("request_id", requestID)
	}
	c.Next()
}

func GenRequestID() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}
