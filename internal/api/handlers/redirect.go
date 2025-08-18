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
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short code"})
        return
    }

    // Get URL from storage
    urlData, err := h.shortenerService.GetURL(c.Request.Context(), shortCode)
    if err != nil {
        if err == url.ErrURLNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
        return
    }

    // Redirect to original URL
    c.Redirect(http.StatusMovedPermanently, urlData.Original)
}