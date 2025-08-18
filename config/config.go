package config

import (
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
    Port       string
    BaseURL    string
    CodeLength int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
    // Load .env file if it exists
    godotenv.Load()
    
    // Get base URL from environment or use default
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:3000"
    }
    
    // Get code length from environment or use default
    codeLengthStr := os.Getenv("CODE_LENGTH")
    codeLength := 6 // Default
    if codeLengthStr != "" {
        if val, err := strconv.Atoi(codeLengthStr); err == nil && val > 0 {
            codeLength = val
        }
    }
    
    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    
    return &Config{
        Port:       port,
        BaseURL:    baseURL,
        CodeLength: codeLength,
    }, nil
}