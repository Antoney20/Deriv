package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init db connection
func ConnectDatabase() error {
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "admin")
	dbName := getEnv("DB_NAME", "go_crud")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbName, host, port)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// // Auto migrate your models here
	// if err := DB.AutoMigrate(&model.User{}); err != nil {
	// 	return err
	// }

	// log.Println("Database migrated successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
