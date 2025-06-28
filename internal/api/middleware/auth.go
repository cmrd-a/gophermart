package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := c.GetHeader("Authorization")
		if a == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		ai, _ := strconv.Atoi(a)

		c.Set("userID", int64(ai))

		c.Next()

		// after request

	}
}
