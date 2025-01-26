package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			c.Abort()
			return
		}

		tokenSplit := strings.Split(token, " ")

		if len(tokenSplit) != 2 || strings.ToLower(tokenSplit[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token, expects bearer token"})
			c.Abort()
			return
		}

		fmt.Println("Authorization Token:", token)
		fmt.Println("Token Parts:", tokenSplit)

		userId, role, err := s.tokenController.VerifyToken(tokenSplit[1])
		if err != nil {
			fmt.Println("Token verification failed:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			c.Abort()
			return
		}

		fmt.Println("User ID from token:", userId)
		fmt.Println("User Role from token:", role)

		// Store user_id and role in context
		c.Set("user_id", userId)
		c.Set("role", role)

		c.Next()
	}
}

func RoleBasedMiddleware(server *Server, role string) gin.HandlerFunc {
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
		user, err := server.queries.GetUserByID(context.Background(), userID)
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
