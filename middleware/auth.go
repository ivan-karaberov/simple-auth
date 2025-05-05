package middleware

import (
	"simpleAuth/config"
	"simpleAuth/errors"
	"simpleAuth/services"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.APIError(c, errors.ErrHeaderIsMissing)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.APIError(c, errors.ErrInvalidHeaderFormat)
			c.Abort()
			return
		}

		tokenString := parts[1]

		payload, err := services.ValidateToken(tokenString, cfg.RSAPublicKey)

		if err != nil {
			errors.APIError(c, errors.ErrIncorrectToken)
			c.Abort()
			return
		}

		_, err = services.CheckSessionExists(db, payload.SID)
		if err != nil {
			errors.APIError(c, errors.ErrIncorrectToken)
			c.Abort()
			return
		}

		c.Set("sessionID", payload.SID)
		c.Set("userID", payload.Subject)

		c.Next()
	}
}
