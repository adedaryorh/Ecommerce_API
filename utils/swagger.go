package utils

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

func AddSwaggerRoutes(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	app.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
}
