package model

import (
    "errors"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
    "regexp"
)

type User struct {
    gorm.Model
    Username    string `gorm:"unique;not null"`
    PhoneNumber string `gorm:"unique;not null"`
    Password    string `gorm:"not null"`
}

// validation--- registration
func (u *User) Validate(db *gorm.DB) error {
    if u.Username == "" || u.PhoneNumber == "" || u.Password == "" {
        return errors.New("all fields are required")
    }

    var existingUser User
    if err := db.Where("username = ?", u.Username).First(&existingUser).Error; err == nil {
        return errors.New("username is already taken")
    }

    if err := db.Where("phone_number = ?", u.PhoneNumber).First(&existingUser).Error; err == nil {
        return errors.New("phone number is already registered")
    }

    if err := validatePassword(u.Password); err != nil {
        return err
    }

    return nil
}

// pass validation
func validatePassword(password string) error {
    if len(password) < 5 {
        return errors.New("password must be at least 5 characters long")
    }

    if isNumeric(password) {
        return errors.New("password cannot be entirely numeric")
    }

    if isCommonPassword(password) {
        return errors.New("password is too common")
    }

    if !containsAlphanumeric(password) {
        return errors.New("password must contain both letters and numbers")
    }

    return nil
}


func isNumeric(s string) bool {
    for _, c := range s {
        if c < '0' || c > '9' {
            return false
        }
    }
    return true
}

// fot common passwords
func isCommonPassword(password string) bool {
    commonPasswords := []string{"123456", "password", "123456789", "qwerty", "abc123"}
    for _, common := range commonPasswords {
        if password == common {
            return true
        }
    }
    return false
}

//  alphanumeric-- strong psswd
func containsAlphanumeric(s string) bool {
    hasLetter := false
    hasNumber := false

    for _, c := range s {
        if regexp.MustCompile(`[a-zA-Z]`).MatchString(string(c)) {
            hasLetter = true
        }
        if regexp.MustCompile(`[0-9]`).MatchString(string(c)) {
            hasNumber = true
        }
    }
    return hasLetter && hasNumber
}

//  finally hash passwd
func (u *User) HashPassword() {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        panic(err) 
    }
    u.Password = string(hashedPassword)
}


//  login -- check passwd
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
