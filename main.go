package main

import (
	"fmt"
	"github.com/adedaryorh/ecommerceapi/api"
)

func main() {
	fmt.Println("Hello welcome to world of ecommerce ")
	server := api.NewServer(".")
	server.Start(8000)
}
