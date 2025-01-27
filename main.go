package main

import (
	"fmt"
	api "github.com/adedaryorh/ecommerceapi/api"
	_ "github.com/adedaryorh/ecommerceapi/docs"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @title Ecommerca Backend Application
// @version 1.0
// @description This is my first version API for an ecommerce simple model.
// @BasePath /
func main() {
	fmt.Println("Hello welcome to world of ecommerce ")
	server := api.NewServer(".")
	server.Start(8000)
}
