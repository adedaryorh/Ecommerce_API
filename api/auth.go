package api_errors

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/adedaryorh/ecommerceapi/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Auth struct {
	server *Server
}

func (a Auth) router(server *Server) {
	a.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("login", a.login)
	serverGroup.POST("register", a.register)
}

// @Summary User Login
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param user body UserParams true "Login Credentials"
// @Success 200 {object} map[string]string "Token response"
// @Failure 400 {object} api_errors.ApiError "Bad Request"
// @Failure 500 {object} api_errors.ApiError "Internal Server Error"
// @Router /auth/login [post]
func (a Auth) login(c *gin.Context) {
	user := new(UserParams)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := a.server.queries.GetUserByEmail(context.Background(), user.Email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or pass"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := utils.VerifyPassword(user.Password, dbUser.HashedPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or pass"})
		return
	}
	token, err := a.server.tokenController.CreateToken(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary User Registration
// @Description Register a new user (admin registration requires admin privileges)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body UserParams true "Registration Details"
// @Success 201 {object} UserResponse "Successful registration"
// @Failure 400 {object} api_errors.ApiError "Bad Request"
// @Failure 401 {object} api_errors.ApiError "Unauthorized"
// @Failure 403 {object} api_errors.ApiError "Forbidden"
// @Failure 500 {object} api_errors.ApiError "Internal Server Error"
// @Router /auth/register [post]
func (a *Auth) register(c *gin.Context) {
	var user UserParams

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate email and username are not empty
	if user.Email == "" || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Username are required"})
		return
	}
	// Ensure that the role is valid (either 'admin' or 'user')
	if user.Role != "admin" && user.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	// If role is 'admin', ensure that the requesting user is an admin (this check should be done for logged-in users)
	if user.Role == "admin" {
		// Fetch user from context (logged-in user)
		value, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, ok := value.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Fetch the user from DB to verify that the requesting user is an admin
		requestingUser, err := a.server.queries.GetUserByID(context.Background(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if requestingUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create admin users"})
			return
		}
	}

	hashedPassword, err := utils.GenerateHashedPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
		Username:       user.Username,
		Role:           user.Role,
	}
	newUser, err := a.server.queries.CreateUser(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Constraint {
			case "users_email_key":
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
				return
			case "users_username_key":
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{}.toUserResponse(&newUser))
}
