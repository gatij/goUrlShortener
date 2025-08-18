package utils

import (
    "strings"

    gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
    // DefaultShortCodeLength is the default length for generated short codes
    DefaultShortCodeLength = 6
    
    // Default alphabet for URL-safe short codes
    // Excludes similar looking characters like 1, l, I, 0, O
    DefaultAlphabet = "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

// GenerateShortCode generates a random short code for URLs using nanoid
func GenerateShortCode(length int) (string, error) {
    if length <= 0 {
        length = DefaultShortCodeLength
    }
    
    // Generate short code using nanoid with our custom alphabet
    return gonanoid.Generate(DefaultAlphabet, length)
}

// GenerateShortURL creates a full shortened URL given a base URL and a short code
func GenerateShortURL(baseURL, shortCode string) string {
    // Ensure baseURL doesn't end with a slash
    baseURL = strings.TrimSuffix(baseURL, "/")
    
    return baseURL + "/" + shortCode
}