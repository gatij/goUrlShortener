package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/gatij/goUrlShortener/config"
    "github.com/gatij/goUrlShortener/internal/api"
    "github.com/gatij/goUrlShortener/internal/service"
    "github.com/gatij/goUrlShortener/internal/storage/metrics"
    "github.com/gatij/goUrlShortener/internal/storage/url"
)

func main() {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found or error loading it. Using environment variables.")
    }

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Set Gin mode based on environment
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Initialize storage
    urlStore := url.NewMemoryStorage()
    metricsStore := metrics.NewMemoryStorage()

    // Initialize services
    metricsService := service.NewMetricsService(metricsStore)
    shortenerConfig := service.ShortenerConfig{
        BaseURL:    cfg.BaseURL,
        CodeLength: cfg.CodeLength,
    }
    shortenerService := service.NewShortenerService(urlStore, metricsService, shortenerConfig)

    // Setup router
    router := api.SetupRouter(shortenerService, metricsService)

    // Configure server
    server := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in a goroutine
    go func() {
        log.Printf("Starting server on port %s...", cfg.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shut down the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Create context with timeout for shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited properly")
}