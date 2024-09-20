package controller

import (
	"net/http"

	"example.com/myapi/config"
	"example.com/myapi/model"
	"github.com/gin-gonic/gin"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var user model.User

	// Bind JSON for user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data now"})
		return
	}

	// Validate the user model
	if err := user.Validate(config.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	user.HashPassword()

	// Save the user to the database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Identifier string `json:"identifier"` // user can use phone number or username to login
		Password   string `json:"password"`
	}

	// Bind JSON for login
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	var user model.User

	// Find user by username or phone number
	if err := config.DB.Where("username = ? OR phone_number = ?", loginData.Identifier, loginData.Identifier).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate password
	if !user.CheckPassword(loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome! "})
}
