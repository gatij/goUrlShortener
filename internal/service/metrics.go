package service

import (
    "context"
    "errors"

    "github.com/gatij/goUrlShortener/internal/model"
    "github.com/gatij/goUrlShortener/internal/storage/metrics"
)

// MetricsService handles URL metrics operations
type MetricsService struct {
    metricsStore metrics.Storage // Metrics storage
}

// NewMetricsService creates a new metrics service
func NewMetricsService(metricsStore metrics.Storage) *MetricsService {
    return &MetricsService{
        metricsStore: metricsStore,
    }
}

// GetTopDomains retrieves the top N most shortened domains
func (s *MetricsService) GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error) {
    if limit <= 0 {
        limit = 3 // Default to top 3 domains
    }
    
    return s.metricsStore.GetTopDomains(ctx, limit)
}

// IncrementDomainShortenCount increments the shorten count for a domain
func (s *MetricsService) IncrementDomainShortenCount(ctx context.Context, domain string) error {
    // Try to get existing metrics
    metrics, err := s.getDomainMetrics(ctx, domain)
    if err != nil {
        // Create new metrics if not found
        metrics = model.DomainMetrics{
            Domain:       domain,
            ShortenCount: 1, // Initialize with 1 for new domains
        }
    } else {
        // Increment counter for existing domains
        metrics.ShortenCount++
    }
    
    // Save updated metrics
    return s.metricsStore.SaveDomainMetrics(ctx, metrics)
}
