package jwt

import (
	"core/config"
	"core/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Error:  "missing-access-token",
				Detail: "Missing access token",
			})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := ValidateAccessToken(cfg, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Error:  "invalid-access-token",
				Detail: "Invalid access token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)
		c.Next()
	}
}
