package service

import (
    "context"
    "errors"
    "time"

    "github.com/gatij/goUrlShortener/internal/model"
    urlStorage "github.com/gatij/goUrlShortener/internal/storage/url"
    "github.com/gatij/goUrlShortener/pkg/utils"
)

var (
    // ErrInvalidURL is returned when the URL is invalid
    ErrInvalidURL = errors.New("invalid URL format")
)

// ShortenerConfig contains configuration for the URL shortener service
type ShortenerConfig struct {
    BaseURL    string // Base URL for generating short links (e.g., "https://short.io")
    CodeLength int    // Length of generated short codes
}

// ShortenerService handles URL shortening operations
type ShortenerService struct {
    urlStore      urlStorage.Storage
    metricsService *MetricsService
    config        ShortenerConfig
}

// NewShortenerService creates a new shortener service
func NewShortenerService(
    urlStore urlStorage.Storage, 
    metricsService *MetricsService, 
    config ShortenerConfig,
) *ShortenerService {
    return &ShortenerService{
        urlStore:      urlStore,
        metricsService: metricsService,
        config:        config,
    }
}

// CreateShortURL creates a new shortened URL
func (s *ShortenerService) CreateShortURL(ctx context.Context, originalURL string) (model.URL, error) {
    // Validate URL
    urlInfo, err := utils.ProcessURL(originalURL, true)
    if err != nil {
        return model.URL{}, ErrInvalidURL
    }
    
    // Use normalized URL with HTTPS
    normalizedURL := urlInfo.NormalizedURL
    
    // Check if URL already exists in storage
    existingURL, err := s.urlStore.GetByOriginalURL(ctx, normalizedURL)
    if err == nil {
        // URL already exists, return it
        // No need to update metrics as it's not a new shortening
        return existingURL, nil
    } else if err != urlStorage.ErrURLNotFound {
        // Unexpected error occurred
        return model.URL{}, err
    }
    
    // URL doesn't exist, create a new short code
    shortCode, err := utils.GenerateShortCode(s.config.CodeLength)
    if err != nil {
        return model.URL{}, err
    }
    
    // Create URL record - using same value for ID and ShortCode
    url := model.URL{
        ID:        shortCode, // Using shortCode as ID
        ShortCode: shortCode,
        Original:  normalizedURL,
        CreatedAt: time.Now(),
    }
    
    // Save URL
    if err := s.urlStore.Save(ctx, url); err != nil {
        return model.URL{}, err
    }
    
    // Extract domain and update metrics asynchronously
    // Only increment metrics for new URLs
    domain := urlInfo.Domain
    go s.metricsService.IncrementDomainShortenCount(ctx, domain)
    
    return url, nil
}

// GetURL retrieves a URL by its short code
func (s *ShortenerService) GetURL(ctx context.Context, shortCode string) (model.URL, error) {
    return s.urlStore.GetByShortCode(ctx, shortCode)
}

// GenerateShortURL creates the full shortened URL given a short code
func (s *ShortenerService) GenerateShortURL(shortCode string) string {
    return utils.GenerateShortURL(s.config.BaseURL, shortCode)
}