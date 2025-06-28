package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := c.GetHeader("Authorization")
		if a == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("userLogin", a)

		// before request

		c.Next()

		// after request

	}
}
