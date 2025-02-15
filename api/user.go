package api_errors

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
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
	serverGroup.DELETE("/:id", u.deleteUser)
	serverGroup.PUT("/:id/password", u.updateUserPassword)
}

// @Summary Delete User
// @Description Delete a user account
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} api_errors.ApiError
// @Failure 403 {object} api_errors.ApiError
// @Failure 404 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /users/{id} [delete]
func (u *User) deleteUser(c *gin.Context) {
	// Check if the logged-in user is an admin
	value, exists := c.Get("user_role")
	if !exists || value.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only administrators can delete users"})
		return
	}

	id := c.Param("id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = u.server.queries.DeleteUser(context.Background(), userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// @Summary Get User By ID
// @Description Retrieve a user by their ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} api_errors.ApiError
// @Failure 404 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /users/{id} [get]
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

// @Summary List Users
// @Description Retrieve a list of users
// @Tags Users
// @Produce json
// @Success 200 {array} UserResponse
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /users [get]
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

// @Summary Get Logged-In User
// @Description Retrieve details of the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /users/me [get]
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

// @Summary Update User Password
// @Description Update a user's password
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param password body UpdatePasswordRequest true "Password update parameters"
// @Success 200 {object} UserResponse
// @Failure 400 {object} api_errors.ApiError
// @Failure 401 {object} api_errors.ApiError
// @Failure 403 {object} api_errors.ApiError
// @Failure 404 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /users/{id}/password [put]
func (u *User) updateUserPassword(c *gin.Context) {
	// Get the target user ID from the URL
	id := c.Param("id")
	targetUserID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the logged-in user's ID
	loggedInUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if the logged-in user is updating their own password or is an admin
	userRole, _ := c.Get("user_role")
	if loggedInUserID.(int64) != targetUserID && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's password"})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user to verify current password
	user, err := u.server.queries.GetUserByID(context.Background(), targetUserID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify current password (skip for admin users updating others' passwords)
	if loggedInUserID.(int64) == targetUserID {
		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.CurrentPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update the password
	arg := db.UpdateUserPasswordParams{
		ID:             targetUserID,
		HashedPassword: string(hashedPassword),
		UpdatedAt:      time.Now(),
	}

	updatedUser, err := u.server.queries.UpdateUserPassword(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&updatedUser))
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
