package api

import (
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

		userId, err := s.tokenController.VerifyToken(tokenSplit[1])
		if err != nil {
			fmt.Println("Token verification failed:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		fmt.Println("User ID from token:", userId)

		c.Set("user_id", userId)
		c.Next()
	}
}
