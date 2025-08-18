package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gatij/goUrlShortener/internal/service"
)

// URLRequest represents the request to create a shortened URL
type URLRequest struct {
    URL string `json:"url" binding:"required"`
}

// URLResponse represents the response with the shortened URL
type URLResponse struct {
    ShortCode  string `json:"short_code"`
    ShortURL   string `json:"short_url"`
    OriginalURL string `json:"original_url"`
}

// ShortenerHandler handles URL shortening endpoints
type ShortenerHandler struct {
    shortenerService *service.ShortenerService
}

// NewShortenerHandler creates a new shortener handler
func NewShortenerHandler(shortenerService *service.ShortenerService) *ShortenerHandler {
    return &ShortenerHandler{
        shortenerService: shortenerService,
    }
}

// CreateShortURL handles requests to create a shortened URL
func (h *ShortenerHandler) CreateShortURL(c *gin.Context) {
    var req URLRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    url, err := h.shortenerService.CreateShortURL(c.Request.Context(), req.URL)
    if err != nil {
        if err == service.ErrInvalidURL {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
        return
    }

    shortURL := h.shortenerService.GenerateShortURL(url.ShortCode)
    
    c.JSON(http.StatusCreated, URLResponse{
        ShortCode:   url.ShortCode,
        ShortURL:    shortURL,
        OriginalURL: url.Original,
    })
}