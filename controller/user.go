package controller

import (
	"net/http"
	"strconv"

	"example.com/myapi/config"
	"example.com/myapi/model"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data now"})
		return
	}

	if err := model.ValidatePhoneNumber(user.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser model.User
	if err := config.DB.Where("phone_number = ?", user.PhoneNumber).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is already registered"})
		return
	}
	
	if err := user.Validate(config.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.HashPassword()

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

//  handles user login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	// Bind JSON for login
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	var user model.User

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

func CreateProfile(c *gin.Context) {
	var profile model.Profile

	// Bind JSON for profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, profile.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existingProfile model.Profile
	if err := config.DB.First(&existingProfile, "user_id = ?", profile.UserID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Profile already exists. Do you want to update it?",
			"existingProfile": existingProfile,
		})
		return
	}


	profile.UserID = user.ID

	if err := config.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully!"})
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


	var user model.User
	if err := config.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	profile.UserID = user.ID 
	if err := config.DB.Model(&user.Profile).Updates(profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully!"})
}


func GetAllUsers(c *gin.Context) {
	var users []model.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

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
// DeleteUser
func DeleteUser(c *gin.Context) {
	userID := c.Param("userID")

	var user model.User
	if err := config.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Select("Profile").Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
}