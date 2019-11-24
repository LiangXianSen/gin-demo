package middleware

import (
	"github.com/gin-gonic/gin"
)

func AKRequired(c *gin.Context) {
	c.Next()
}

func LoginRequired(c *gin.Context) {
	c.Next()
}
