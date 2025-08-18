package metrics

import (
	"context"
	"testing"

	"github.com/gatij/goUrlShortener/internal/model"
)

func TestMemoryStorage_SaveDomainMetrics(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create test metrics
	metrics := model.DomainMetrics{
		Domain:      "example.com",
		ShortenCount: 5,
	}

	// Test saving new metrics
	err := storage.SaveDomainMetrics(ctx, metrics)
	if err != nil {
		t.Errorf("Failed to save domain metrics: %v", err)
	}

	// Test updating existing metrics
	updatedMetrics := model.DomainMetrics{
		Domain:      "example.com",
		ShortenCount: 10,
	}
	err = storage.SaveDomainMetrics(ctx, updatedMetrics)
	if err != nil {
		t.Errorf("Failed to update domain metrics: %v", err)
	}

	// Verify metrics were updated
	retrievedMetrics, exists, err := storage.GetDomainMetrics(ctx, "example.com")
	if err != nil {
		t.Errorf("Failed to get domain metrics: %v", err)
	}
	if !exists {
		t.Errorf("Domain metrics should exist but doesn't")
	}
	if retrievedMetrics.ShortenCount != 10 {
		t.Errorf("Expected ShortenCount 10 but got %d", retrievedMetrics.ShortenCount)
	}
}

func TestMemoryStorage_GetDomainMetrics(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create and save test metrics
	metrics := model.DomainMetrics{
		Domain:      "example.com",
		ShortenCount: 5,
	}

	err := storage.SaveDomainMetrics(ctx, metrics)
	if err != nil {
		t.Fatalf("Failed to save domain metrics: %v", err)
	}

	// Test getting existing metrics
	retrievedMetrics, exists, err := storage.GetDomainMetrics(ctx, "example.com")
	if err != nil {
		t.Errorf("Failed to get domain metrics: %v", err)
	}
	if !exists {
		t.Errorf("Domain metrics should exist but doesn't")
	}
	if retrievedMetrics.Domain != metrics.Domain || retrievedMetrics.ShortenCount != metrics.ShortenCount {
		t.Errorf("Expected metrics %+v but got %+v", metrics, retrievedMetrics)
	}

	// Test getting non-existent metrics
	_, exists, err = storage.GetDomainMetrics(ctx, "nonexistent.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if exists {
		t.Errorf("Domain metrics shouldn't exist but does")
	}
}

func TestMemoryStorage_GetTopDomains(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create and save test metrics for multiple domains
	domains := []model.DomainMetrics{
		{Domain: "example.com", ShortenCount: 5},
		{Domain: "github.com", ShortenCount: 10},
		{Domain: "golang.org", ShortenCount: 3},
		{Domain: "google.com", ShortenCount: 7},
	}

	for _, domain := range domains {
		err := storage.SaveDomainMetrics(ctx, domain)
		if err != nil {
			t.Fatalf("Failed to save domain metrics: %v", err)
		}
	}

	// Test getting top 2 domains
	topDomains, err := storage.GetTopDomains(ctx, 2)
	if err != nil {
		t.Errorf("Failed to get top domains: %v", err)
	}
	if len(topDomains) != 2 {
		t.Errorf("Expected 2 domains but got %d", len(topDomains))
	}
	if topDomains[0].Domain != "github.com" || topDomains[0].ShortenCount != 10 {
		t.Errorf("Expected first domain to be github.com with count 10, but got %s with count %d", 
			topDomains[0].Domain, topDomains[0].ShortenCount)
	}
	if topDomains[1].Domain != "google.com" || topDomains[1].ShortenCount != 7 {
		t.Errorf("Expected second domain to be google.com with count 7, but got %s with count %d", 
			topDomains[1].Domain, topDomains[1].ShortenCount)
	}

	// Test getting all domains (limit > number of domains)
	allDomains, err := storage.GetTopDomains(ctx, 10)
	if err != nil {
		t.Errorf("Failed to get all domains: %v", err)
	}
	if len(allDomains) != 4 {
		t.Errorf("Expected 4 domains but got %d", len(allDomains))
	}

	// Test with zero limit (should use default)
	defaultLimitDomains, err := storage.GetTopDomains(ctx, 0)
	if err != nil {
		t.Errorf("Failed to get domains with default limit: %v", err)
	}
	if len(defaultLimitDomains) == 0 {
		t.Errorf("Expected some domains with default limit but got none")
	}
}
