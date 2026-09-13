package middleware

import (
	"net/http"
	"strings"

	"movie-booking/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// middleware verifies the jwt sent by the client before allowing access to protected routes
func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// read authorization header
		authHeader := c.GetHeader("Authorization")

		// check if header is present and starts with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid authorization header",
			})
			c.Abort()
			return

		}
		// extract token string by trimming "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// extract claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// store user id in gin context
		c.Set("user_id", uint(claims["user_id"].(float64)))

		// store user role in gin context
		c.Set("role", claims["role"].(string))

		// continue to next handler
		c.Next()

	}
}
