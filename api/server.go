package api_errors

import (
	"database/sql"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/adedaryorh/ecommerceapi/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	_ "github.com/go-playground/validator/v10"
	"github.com/golodash/galidator"
	_ "github.com/lib/pq"
)

type Server struct {
	queries         *db.Queries
	router          *gin.Engine
	config          *utils.Config
	tokenController *utils.JWTToken
}

var gValid = galidator.New().CustomMessages(
	galidator.Messages{
		"required": "this field is required",
	},
)

func myCorsHandler() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	return cors.New(config)
}

func NewServer(envPath string) *Server {
	config, err := utils.LoadConfig(envPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}
	conn, err := sql.Open(config.DBdriver, config.DB_source_live)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to DB: %v", err))
	}
	tokenController := utils.NewJWTToken(config)

	q := db.New(conn)
	g := gin.Default()
	g.Use(myCorsHandler())

	return &Server{
		queries:         q,
		router:          g,
		config:          config,
		tokenController: tokenController,
	}
}

func (s *Server) initializeRoutes() {
	router := s.router
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Public routes (accessible by any authenticated user)
	// Public routes (accessible by any authenticated user)
	// @Summary Create Order
	// @Description Create a new order (authenticated user)
	// @Tags Orders
	// @Accept json
	// @Produce json
	// @Param order body OrderParams true "Order Details"
	// @Success 201 {object} OrderResponse
	// @Failure 400 {object} api_errors.ApiError
	// @Failure 500 {object} api_errors.ApiError
	// @Security BearerAuth
	// @Router /orders [post]
	router.POST("/orders", s.AuthenticatedMiddleware(), s.CreateOrder)
	// @Summary List User Orders
	// @Description Retrieve a list of orders placed by the authenticated user
	// @Tags Orders
	// @Produce json
	// @Success 200 {array} OrderResponse
	// @Failure 500 {object} api_errors.ApiError
	// @Security BearerAuth
	// @Router /orders [get]
	router.GET("/orders", s.AuthenticatedMiddleware(), s.ListUserOrders)

	// Admin routes (only accessible by admins)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(s.AuthenticatedMiddleware(), RoleBasedMiddleware(s, "admin"))
	{
		// @Summary Cancel Order
		// @Description Cancel an order by ID (admin only)
		// @Tags Orders
		// @Param id path string true "Order ID"
		// @Success 200 {object}
		// @Failure 400 {object} api_errors.ApiError
		// @Failure 404 {object} api_errors.ApiError
		// @Failure 500 {object} api_errors.ApiError
		// @Security BearerAuth
		// @Router /admin/orders/{id}/cancel [post]
		adminRoutes.POST("/orders/:id/cancel", s.CancelOrder)
		// @Summary Update Order Status
		// @Description Update the status of an order (admin only)
		// @Tags Orders
		// @Param id path string true "Order ID"
		// @Param status body string true "New Order Status"
		// @Success 200 {object}
		// @Failure 400 {object} api_errors.ApiError
		// @Failure 404 {object} api_errors.ApiError
		// @Failure 500 {object} api_errors.ApiError
		// @Security BearerAuth
		// @Router /admin/orders/{id}/status [patch]
		adminRoutes.PATCH("/orders/:id/status", s.UpdateOrderStatus)
	}

	// Assign router to the server instance
	s.router = router
}

func (s *Server) Start(port int) {
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Ecommerce API"})
	})

	User{}.router(s)
	Auth{}.router(s)
	(&Product{}).router(s)
	s.initializeRoutes()

	s.router.Run(fmt.Sprintf(":%v", port))
}
