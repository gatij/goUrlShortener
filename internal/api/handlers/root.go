package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// RootHandler handles requests to the root endpoint
func RootHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "service": "URL Shortener",
        "version": "1.0.0",
        "description": "A simple, scalable URL shortener service written in Go",
        "endpoints": gin.H{
            "create_short_url": "POST /api/v1/urls",
            "get_top_domains": "GET /api/v1/metrics/domains",
            "redirect": "GET /{shortCode}",
            "health": "GET /health",
        },
        "usage_example": "POST /api/v1/urls with {\"url\": \"https://example.com/long/url\"}",
        "source_code": "https://github.com/gatij/goUrlShortener",
    })
}