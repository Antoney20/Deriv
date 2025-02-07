package main

import (
	"log"

	"example.com/myapi/config"
	"example.com/myapi/controller"
	"example.com/myapi/model" 
	"example.com/myapi/middleware" 
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err := config.DB.AutoMigrate(&model.User{},  &model.Profile{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Println("Migrations successful")

	router := gin.Default()
    router.Use(middleware.LoggerMiddleware())
	

	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/profile/:id", controller.CreateProfile)
	router.PUT("/profile/:id", controller.UpdateProfile)
	router.GET("/users", controller.GetAllUsers) 
	router.GET("/users/:id", controller.FetchUserByID)
	router.DELETE("/users/:userID", controller.DeleteUser)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
