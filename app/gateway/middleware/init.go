package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// catch Handler (or middleware) panic to prevent from crash
func ErrorMiddleware() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// map[string]interface{}{"key": value}
				c.JSON(200, gin.H{
					"code": 404,
					"msg":  fmt.Sprintf("%s", r),
				})
				c.Abort()
			}
		}()
		c.Next()
	}

	return handler
}
