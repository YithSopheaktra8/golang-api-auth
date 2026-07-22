package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})

			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})

			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// prevent algorithm attack
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

					return nil, jwt.ErrSignatureInvalid

				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})

			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid claims",
			})

			c.Abort()
			return
		}

		// IMPORTANT
		// only access token can call APIs

		tokenType, ok := claims["type"].(string)

		if !ok || tokenType != "access" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token type",
			})

			c.Abort()
			return
		}

		c.Set(
			"user_id",
			claims["user_id"],
		)

		c.Set(
			"email",
			claims["email"],
		)

		c.Next()

	}

}
