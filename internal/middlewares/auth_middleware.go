package middlewares

import (
	"net/http"
	"github.com/KasiditR/forviz-backend-api-test/internal/utils"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No Authorization Header"})
			c.Abort()
			return
		}

		claims, msg := utils.ValidateToken(clientToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": msg})
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("user_name", claims.UserName)
		c.Next()
	}
}
