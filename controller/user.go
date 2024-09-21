package controller

import (
	"net/http"
	"strconv"

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

	// Validate phone number
	if err := model.ValidatePhoneNumber(user.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for existing user by phone number
	var existingUser model.User
	if err := config.DB.Where("phone_number = ?", user.PhoneNumber).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is already registered"})
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


// UpdateProfile handles updating user profiles
func UpdateProfile(c *gin.Context) {
	var profile model.Profile
	userID := c.Param("userID") 

	// Bind JSON for profile update
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Find the user by ID
	var user model.User
	if err := config.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update the profile
	profile.UserID = user.ID // Associate the profile with the user
	if err := config.DB.Model(&user.Profile).Updates(profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully!"})
}


func GetAllUsers(c *gin.Context) {
	var users []model.User

	// Fetch all users with profiles
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// retrieves user by ID 
func FetchUserByID(c *gin.Context) {
	id := c.Param("id") //get user id


	userID, err := strconv.ParseUint(id, 10, 32) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
// DeleteUser handles deleting a user by ID
func DeleteUser(c *gin.Context) {
	userID := c.Param("userID")

	// Find the user by ID
	var user model.User
	if err := config.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user and associated profile
	if err := config.DB.Select("Profile").Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
}