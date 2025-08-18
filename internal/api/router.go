package api

import (
    "github.com/gin-gonic/gin"
    "github.com/gatij/goUrlShortener/internal/api/handlers"
    "github.com/gatij/goUrlShortener/internal/api/middleware"
    "github.com/gatij/goUrlShortener/internal/service"
)

// SetupRouter configures the API routes
func SetupRouter(
    shortenerService *service.ShortenerService, 
    metricsService *service.MetricsService,
) *gin.Engine {
    // Create router with default middleware
    router := gin.Default()

    // Add custom logging middleware
    router.Use(middleware.Logger())

    // Create handlers
    shortenerHandler := handlers.NewShortenerHandler(shortenerService)
    redirectHandler := handlers.NewRedirectHandler(shortenerService)
    metricsHandler := handlers.NewMetricsHandler(metricsService)

    // API routes
    api := router.Group("/api/v1")
    {
        // URL shortening endpoint
        api.POST("/urls", shortenerHandler.CreateShortURL)
        
        // Metrics endpoint
        api.GET("/metrics/domains", metricsHandler.GetTopDomains)
    }

    // Redirect route - must be last to catch all other paths
    router.GET("/:shortCode", redirectHandler.RedirectToOriginal)

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.String(200, "OK")
    })

    return router
}