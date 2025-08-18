package service

import (
	"context"
	"testing"

	"github.com/gatij/goUrlShortener/internal/model"
)

// MockMetricsStorage is a mock implementation of the metrics storage interface
type MockMetricsStorage struct {
	domains map[string]model.DomainMetrics
}

func NewMockMetricsStorage() *MockMetricsStorage {
	return &MockMetricsStorage{
		domains: make(map[string]model.DomainMetrics),
	}
}

func (m *MockMetricsStorage) SaveDomainMetrics(ctx context.Context, metrics model.DomainMetrics) error {
	m.domains[metrics.Domain] = metrics
	return nil
}

func (m *MockMetricsStorage) GetDomainMetrics(ctx context.Context, domain string) (model.DomainMetrics, bool, error) {
	metrics, exists := m.domains[domain]
	return metrics, exists, nil
}

func (m *MockMetricsStorage) GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error) {
	domains := make([]model.DomainMetrics, 0, len(m.domains))
	for _, metrics := range m.domains {
		domains = append(domains, metrics)
	}
	
	// Simple sort by shorten count (not efficient but works for tests)
	for i := 0; i < len(domains); i++ {
		for j := i + 1; j < len(domains); j++ {
			if domains[i].ShortenCount < domains[j].ShortenCount {
				domains[i], domains[j] = domains[j], domains[i]
			}
		}
	}
	
	if limit > 0 && limit < len(domains) {
		domains = domains[:limit]
	}
	
	return domains, nil
}

func TestMetricsService_GetTopDomains(t *testing.T) {
	// Set up mock
	metricsStorage := NewMockMetricsStorage()
	
	// Add some test domains
	domains := []model.DomainMetrics{
		{Domain: "example.com", ShortenCount: 5},
		{Domain: "github.com", ShortenCount: 10},
		{Domain: "golang.org", ShortenCount: 3},
	}
	
	ctx := context.Background()
	for _, domain := range domains {
		err := metricsStorage.SaveDomainMetrics(ctx, domain)
		if err != nil {
			t.Fatalf("Failed to save domain metrics: %v", err)
		}
	}
	
	// Create service
	service := NewMetricsService(metricsStorage)
	
	// Test getting top 2 domains
	topDomains, err := service.GetTopDomains(ctx, 2)
	if err != nil {
		t.Fatalf("Failed to get top domains: %v", err)
	}
	
	// Verify results
	if len(topDomains) != 2 {
		t.Errorf("Expected 2 domains but got %d", len(topDomains))
	}
	if topDomains[0].Domain != "github.com" || topDomains[0].ShortenCount != 10 {
		t.Errorf("Expected first domain to be github.com with count 10, but got %s with count %d", 
			topDomains[0].Domain, topDomains[0].ShortenCount)
	}
	if topDomains[1].Domain != "example.com" || topDomains[1].ShortenCount != 5 {
		t.Errorf("Expected second domain to be example.com with count 5, but got %s with count %d", 
			topDomains[1].Domain, topDomains[1].ShortenCount)
	}
	
	// Test with default limit
	defaultTopDomains, err := service.GetTopDomains(ctx, 0)
	if err != nil {
		t.Fatalf("Failed to get top domains with default limit: %v", err)
	}
	
	// Verify results (default should be 3)
	if len(defaultTopDomains) != 3 {
		t.Errorf("Expected 3 domains with default limit but got %d", len(defaultTopDomains))
	}
}

func TestMetricsService_IncrementDomainShortenCount(t *testing.T) {
	// Set up mock
	metricsStorage := NewMockMetricsStorage()
	
	// Create service
	service := NewMetricsService(metricsStorage)
	ctx := context.Background()
	
	// Test incrementing a new domain
	err := service.IncrementDomainShortenCount(ctx, "example.com")
	if err != nil {
		t.Fatalf("Failed to increment domain count: %v", err)
	}
	
	// Verify domain was added with count 1
	metrics, exists, err := metricsStorage.GetDomainMetrics(ctx, "example.com")
	if err != nil {
		t.Fatalf("Failed to get domain metrics: %v", err)
	}
	if !exists {
		t.Fatalf("Domain metrics should exist but doesn't")
	}
	if metrics.ShortenCount != 1 {
		t.Errorf("Expected shorten count to be 1 but got %d", metrics.ShortenCount)
	}
	
	// Test incrementing an existing domain
	err = service.IncrementDomainShortenCount(ctx, "example.com")
	if err != nil {
		t.Fatalf("Failed to increment domain count: %v", err)
	}
	
	// Verify count was incremented
	metrics, exists, err = metricsStorage.GetDomainMetrics(ctx, "example.com")
	if err != nil {
		t.Fatalf("Failed to get domain metrics: %v", err)
	}
	if !exists {
		t.Fatalf("Domain metrics should exist but doesn't")
	}
	if metrics.ShortenCount != 2 {
		t.Errorf("Expected shorten count to be 2 but got %d", metrics.ShortenCount)
	}
}
