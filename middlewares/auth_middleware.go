package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("qC2dACu+zgx94ALrfmCTESkxoqfCG4ItCWknbz+XmfTfWNDvFeuYKOGXgAKqSq+7Bdu8jXrxZfpWwE0K0jPLHw==") // Replace with the same secret key used for token creation.

// JWTAuthMiddleware is a middleware to validate the JWT token.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Check if the Authorization header is provided and starts with "Bearer ".
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header invalid"})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Token is valid, continue to the next middleware or route handler.
		c.Next()
	}
}
