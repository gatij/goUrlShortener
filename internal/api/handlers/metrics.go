package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/gatij/goUrlShortener/internal/service"
)

// MetricsHandler handles metrics endpoints
type MetricsHandler struct {
    metricsService *service.MetricsService
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(metricsService *service.MetricsService) *MetricsHandler {
    return &MetricsHandler{
        metricsService: metricsService,
    }
}

// GetTopDomains returns the top N most shortened domains
func (h *MetricsHandler) GetTopDomains(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "3")
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 3 // Default to top 3 if invalid
    }

    domains, err := h.metricsService.GetTopDomains(c.Request.Context(), limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top domains"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "top_domains": domains,
        "limit":       limit,
    })
}