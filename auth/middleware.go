package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ClaimsContextKey = "claims"

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// fallback: try getting token from query parameter
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			if strings.Contains(c.FullPath(), "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			if strings.Contains(c.FullPath(), "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}

		c.Set(ClaimsContextKey, claims)
		c.Next()
	}
}
