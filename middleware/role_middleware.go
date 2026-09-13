package middleware

import (
	"net/http"

	"movie-booking/models"

	"github.com/gin-gonic/gin"
)

//this middleware always run after authmiddleware
/*because AuthMiddleware is responsible for:

1. validating the JWT
2. extracting user information from the token
3. storing user_id and role inside Gin's Context

admin_middleware simply reads role from gin's context and allows only admins to continue

*/

func AdminMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		//retrieve user's role from gin's conetext
		//this value was stored earlier by authmiddleware we used c.set smthing
		role, exists := c.Get("role")

		//if no role smth is wrong
		//most likely authmiddleware was skipped
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "role not found",
			})
			c.Abort()
			return
		}

		//check if logged in user is admin
		if role != models.RoleAdmin {

			//user is authenticated but doesnt have permission
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})

			//stop executing any remaining handlers
			c.Abort()
			return
		}

		//user is admin
		//continue to the requested route
		c.Next()
	}
}
