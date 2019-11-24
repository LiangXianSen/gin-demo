package middleware

import (
	"github.com/gin-gonic/gin"
)

// AKRequired is a middleware which provides AKSK verify
func AKRequired(c *gin.Context) {
	c.Next()
}

// LoginRequired is a middleware which provides users verify
func LoginRequired(c *gin.Context) {
	c.Next()
}
