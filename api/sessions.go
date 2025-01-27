package api_errors

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/google/uuid"
)

type createSessionRequest struct {
	UserID int64 `json:"user_id"`
}

type sessionResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateSession handles session creation
func (server *Server) CreateSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Generate a unique session token
	token := uuid.New().String()

	// Set session expiry (e.g., 24 hours from now)
	expiresAt := time.Now().UTC().Add(24 * time.Hour)

	arg := db.CreateSessionParams{
		UserID:    req.UserID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	session, err := server.queries.CreateSession(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	// Set session cookie
	c.SetCookie(
		"session_token",
		token,
		int(24*time.Hour.Seconds()), // MaxAge in seconds
		"/",                         // Path
		"",                          // Domain
		true,                        // Secure
		true,                        // HttpOnly
	)

	c.JSON(http.StatusCreated, sessionResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	})
}

// GetSession retrieves session information
func (server *Server) GetSession(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no session found"})
		return
	}

	session, err := server.queries.GetSessionByToken(c, token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get session"})
		return
	}

	// Check if session has expired
	if time.Now().UTC().After(session.ExpiresAt) {
		err := server.queries.DeleteSession(c, token)
		if err != nil {
			// Use Gin's logger for error logging
			c.Error(err) // This will be logged by Gin
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
		return
	}

	c.JSON(http.StatusOK, sessionResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	})
}

// DeleteSession handles session deletion (logout)
func (server *Server) DeleteSession(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no session found"})
		return
	}

	err = server.queries.DeleteSession(c, token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete session"})
		return
	}

	// Remove the session cookie
	c.SetCookie(
		"session_token",
		"",
		-1, // MaxAge < 0 means delete cookie now
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "session deleted successfully",
	})
}

// Add these routes to your server setup
func (server *Server) setupSessionRoutes() {
	sessionGroup := server.router.Group("/api/sessions")
	{
		sessionGroup.POST("", server.CreateSession)
		sessionGroup.GET("", server.GetSession)
		sessionGroup.DELETE("", server.DeleteSession)
	}
}
