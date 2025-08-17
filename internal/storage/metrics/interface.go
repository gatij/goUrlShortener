package metrics

import (
	"context"

	"github.com/gatij/goUrlShortener/internal/model"
)

// Storage defines the interface for metrics storage operations
type Storage interface {
	// SaveDomainMetrics stores metrics for a domain
	SaveDomainMetrics(ctx context.Context, metrics model.DomainMetrics) error

	// GetTopDomains retrieves the top N domains based on shorten count
	GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error)
}