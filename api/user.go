package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/gin-gonic/gin"
)

type User struct {
	server *Server
}

func (u User) router(server *Server) {
	u.server = server

	serverGroup := server.router.Group("/users", u.server.AuthenticatedMiddleware())
	serverGroup.GET("", u.listUsers)
	serverGroup.GET("/:id", u.getUserByID)
	serverGroup.GET("/me", u.getLoggedInUser)
}

func (u *User) getUserByID(c *gin.Context) {
	id := c.Param("id")

	// Convert the ID to an int64
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := u.server.queries.GetUserByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *User) listUsers(c *gin.Context) {
	arg := db.ListUserParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := u.server.queries.ListUser(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUsers := []UserResponse{}

	for _, v := range users {
		n := UserResponse{}.toUserResponse(&v)
		newUsers = append(newUsers, *n)
	}

	c.JSON(http.StatusOK, newUsers)
}

type UserParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

func (u *User) getLoggedInUser(c *gin.Context) {
	value, exists := c.Get("user_id") // Change this to match the key you set
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized to access resources"})
		return
	}
	userId, ok := value.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encountered an err"})
		return
	}
	user, err := u.server.queries.GetUserByID(context.Background(), userId)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or pass"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(user *db.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

}
