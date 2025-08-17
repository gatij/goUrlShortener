package main

import (
	"fmt"
)

const (
    defaultPort = "3000"
)

func main() {
	// Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("Error while loading .env file.")
    }

	// Configure port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = defaultPort
    }

	// Create a new Gin router
    router := gin.Default()

	// Root path handler
    router.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "URL Shortener Service - API endpoints will be available soon")
    })

	// Redirect handler for shortened URLs
    // This captures the shortCode parameter from the URL
    router.GET("/:shortCode", func(c *gin.Context) {
        shortCode := c.Param("shortCode")
        
        // For now, just show what we captured
        // This will be replaced with actual URL lookup and redirect logic
        c.String(http.StatusOK, "Will redirect short code: %s", shortCode)

    })

	// API routes will be added here
    apiV1 := router.Group("/api/v1")
    {
        // Placeholder for metrics endpoint
		// This will eventually return the top 3 hot domains
		apiV1.GET("/metrics", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "Top 3 hot domains will be displayed here",
            })
        })
        
        // Placeholder for URL shortening endpoint
        apiV1.POST("/urls", func(c *gin.Context) {
            c.JSON(http.StatusNotImplemented, gin.H{
                "message": "URL shortening coming soon",
            })
        })
    }
}