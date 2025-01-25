package api

import (
	"database/sql"
	"fmt"
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
	//tokenController = utils.NewJWTToken(config)

	q := db.New(conn)

	//g.Use(myCorsHandler())

	g := gin.Default()

	return &Server{
		queries:         q,
		router:          g,
		config:          config,
		tokenController: utils.NewJWTToken(config),
	}
}

func (s *Server) Start(port int) {
	// Root route
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Ecommerce API"})
	})

	User{}.router(s)
	Auth{}.router(s)

	s.router.Run(fmt.Sprintf(":%v", port))
}
