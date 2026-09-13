package middleware

import (
	"net/http"
	"strings"

	"movie-booking/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Service returns the JWT to the handler; the handler sends it to the client as the HTTP response.

// Middleware verifies the JWT sent by the client before allowing access to protected routes.
func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		//read authorization header
		authHeader := c.GetHeader("Authorization")

		//check if header is present and starts wirh "Bearer "..
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") { //could be empty or not bearer thats ||
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid authorization header",
			})
			c.Abort()
			return

		}
		//if not then it has bearer and it proceeds further to
		//extract  the jwt by removing "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//validate the jwt
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		//extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		//store user id in gin context
		// Take the value stored under the "user_id" key in the JWT claims.
		//Store it in Gin's request context using the key "user_id".
		// Convert the JWT claim to uint before storing it in Gin context.
		c.Set("user_id", uint(claims["user_id"].(float64)))

		//store user role in gin's contexr
		c.Set("role", claims["role"].(string))

		//continue to next handler
		c.Next()

	}
}
