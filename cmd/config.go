package main

import (
	"fmt"
	"os"
	// If you're using godotenv, uncomment the following line
	// "github.com/joho/godotenv"
)

// DBConfig holds the database configuration
type DBConfig struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     string
}

// LoadConfig loads the database configuration from environment variables
func LoadConfig() (*DBConfig, error) {
	// If you're using godotenv, uncomment the following line
	// err := godotenv.Load()
	// if err != nil {
	//     return nil, fmt.Errorf("error loading .env file: %v", err)
	// }

	// Read environment variables
	dbConfig := &DBConfig{
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}

	// Check if all required fields are populated
	if dbConfig.DBName == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.Host == "" || dbConfig.Port == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return dbConfig, nil
}
