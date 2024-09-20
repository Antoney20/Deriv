package config

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() error {
    // Retrieve database connection details from environment variables
    user := getEnv("DB_USER", "postgres")
    password := getEnv("DB_PASSWORD", "admin")
    dbName := getEnv("DB_NAME", "go_crud")
    host := getEnv("DB_HOST", "localhost")
    port := getEnv("DB_PORT", "5432")

    // Construct the Data Source Name (DSN)
    dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
        user, password, dbName, host, port)

    // Connect to the database using gorm
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }

    log.Println("Database connection established")
    return nil
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
