package main

import (
	"log"

	"example.com/myapi/config"
	"example.com/myapi/controller"
	"example.com/myapi/model" 
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// load env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Connect to the database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Auto migrate models
	if err := config.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Println("Migrations successful")

	router := gin.Default()

	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
