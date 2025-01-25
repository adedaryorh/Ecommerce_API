package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleBasedMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, ok := value.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Fetch the user from the database using the userID
		user, err := c.MustGet("server").(*Server).queries.GetUserByID(context.Background(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if the user has the required role
		if user.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
