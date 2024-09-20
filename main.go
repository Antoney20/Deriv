package main

import (
	"log"

	"example.com/myapi/config"
	"example.com/myapi/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	router := gin.Default()

	router.POST("/register", controller.RegisterUser)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
