package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gatij/goUrlShortener/internal/service"
    "github.com/gatij/goUrlShortener/internal/storage/url"
)

// RedirectHandler handles URL redirection
type RedirectHandler struct {
    shortenerService *service.ShortenerService
}

// NewRedirectHandler creates a new redirect handler
func NewRedirectHandler(shortenerService *service.ShortenerService) *RedirectHandler {
    return &RedirectHandler{
        shortenerService: shortenerService,
    }
}

// RedirectToOriginal redirects short URLs to their original destination
func (h *RedirectHandler) RedirectToOriginal(c *gin.Context) {
    shortCode := c.Param("shortCode")
    if shortCode == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Missing short code",
            "message": "The URL is missing a short code. Please use a valid shortened URL.",
            "example": "Try visiting /{shortCode} where shortCode is the code from your shortened URL",
        })
        return
    }

    // Get URL from storage
    urlData, err := h.shortenerService.GetURL(c.Request.Context(), shortCode)
    if err != nil {
        if err == url.ErrURLNotFound {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "URL not found",
                "message": "The shortened URL you're trying to access doesn't exist or has expired.",
                "help": "Please check the URL and try again, or create a new shortened URL at /api/v1/urls.",
                "create_url_endpoint": "/api/v1/urls",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve URL",
            "message": "An internal error occurred while processing your request.",
            "contact": "Please try again later or contact the administrator if the problem persists.",
        })
        return
    }

    // Redirect to original URL
    c.Redirect(http.StatusMovedPermanently, urlData.Original)
}