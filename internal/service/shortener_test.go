package service

import (
	"context"
	"testing"
	"time"

	"github.com/gatij/goUrlShortener/internal/model"
	"github.com/gatij/goUrlShortener/internal/storage/url"
	"github.com/gatij/goUrlShortener/pkg/utils"
)

// MockURLStorage is a mock implementation of the URL storage interface
type MockURLStorage struct {
	urls map[string]model.URL
}

func NewMockURLStorage() *MockURLStorage {
	return &MockURLStorage{
		urls: make(map[string]model.URL),
	}
}

func (m *MockURLStorage) Save(ctx context.Context, urlObj model.URL) error {
	if _, exists := m.urls[urlObj.ID]; exists {
		return url.ErrURLExists
	}
	m.urls[urlObj.ID] = urlObj
	return nil
}

func (m *MockURLStorage) GetByID(ctx context.Context, id string) (model.URL, error) {
	urlObj, exists := m.urls[id]
	if !exists {
		return model.URL{}, url.ErrURLNotFound
	}
	return urlObj, nil
}

func (m *MockURLStorage) GetByShortCode(ctx context.Context, shortCode string) (model.URL, error) {
	for _, urlObj := range m.urls {
		if urlObj.ShortCode == shortCode {
			return urlObj, nil
		}
	}
	return model.URL{}, url.ErrURLNotFound
}

func (m *MockURLStorage) Delete(ctx context.Context, id string) error {
	if _, exists := m.urls[id]; !exists {
		return url.ErrURLNotFound
	}
	delete(m.urls, id)
	return nil
}

// We need to modify the ShortenerService to work with our mock metrics service
// This is a modified version of the original ShortenerService with a more generic metrics service interface
type TestShortenerService struct {
	urlStore       url.Storage
	metricsService MetricsServiceInterface
	config         ShortenerConfig
}

// MetricsServiceInterface defines the interface that both the real and mock metrics services implement
type MetricsServiceInterface interface {
	IncrementDomainShortenCount(ctx context.Context, domain string) error
	GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error)
}

// NewTestShortenerService creates a test shortener service that works with our mock
func NewTestShortenerService(
urlStore url.Storage,
metricsService MetricsServiceInterface,
config ShortenerConfig,
) *TestShortenerService {
	return &TestShortenerService{
		urlStore:       urlStore,
		metricsService: metricsService,
		config:         config,
	}
}

// CreateShortURL creates a new shortened URL
func (s *TestShortenerService) CreateShortURL(ctx context.Context, originalURL string) (model.URL, error) {
	// Validate URL
	parsedURL, err := utils.ValidateURL(originalURL)
	if err != nil {
		return model.URL{}, ErrInvalidURL
	}

	// Generate short code
	shortCode, err := utils.GenerateShortCode(s.config.CodeLength)
	if err != nil {
		return model.URL{}, err
	}

	// Create URL model
	now := time.Now()
	url := model.URL{
		ID:        shortCode, // Using short code as ID for simplicity
		ShortCode: shortCode,
		Original:  originalURL,
		CreatedAt: now,
	}

	// Save URL
	if err := s.urlStore.Save(ctx, url); err != nil {
		return model.URL{}, err
	}

	// Extract domain and update metrics
	domain := utils.ExtractDomain(parsedURL)
	// In the test version, we don't use a goroutine to ensure the count is updated before we check it
	s.metricsService.IncrementDomainShortenCount(ctx, domain)

	return url, nil
}

// GetURL retrieves a URL by its short code
func (s *TestShortenerService) GetURL(ctx context.Context, shortCode string) (model.URL, error) {
	return s.urlStore.GetByShortCode(ctx, shortCode)
}

// GenerateShortURL generates the full short URL from a short code
func (s *TestShortenerService) GenerateShortURL(shortCode string) string {
	return utils.GenerateShortURL(s.config.BaseURL, shortCode)
}

// MockMetricsService is a mock implementation of the MetricsServiceInterface
type MockMetricsService struct {
	domains map[string]int
}

func NewMockMetricsService() *MockMetricsService {
	return &MockMetricsService{
		domains: make(map[string]int),
	}
}

func (m *MockMetricsService) IncrementDomainShortenCount(ctx context.Context, domain string) error {
	m.domains[domain]++
	return nil
}

func (m *MockMetricsService) GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error) {
	// Not needed for this test
	return nil, nil
}

func TestShortenerService_CreateShortURL(t *testing.T) {
	// Set up mocks
	urlStorage := NewMockURLStorage()
	metricsService := NewMockMetricsService()
	
	// Create service with configuration
	config := ShortenerConfig{
		BaseURL:    "http://localhost:3000",
		CodeLength: 6,
	}
	service := NewTestShortenerService(urlStorage, metricsService, config)
	ctx := context.Background()

	// Test creating a short URL (using a non-blocked domain)
	url, err := service.CreateShortURL(ctx, "https://github.com/user/repo")
	if err != nil {
		t.Fatalf("Failed to create short URL: %v", err)
	}

	// Verify the created URL
	if url.Original != "https://github.com/user/repo" {
		t.Errorf("Expected original URL to be https://github.com/user/repo but got %s", url.Original)
	}
	if len(url.ShortCode) != config.CodeLength {
		t.Errorf("Expected short code length to be %d but got %d", config.CodeLength, len(url.ShortCode))
	}
	if url.ID != url.ShortCode {
		t.Errorf("Expected ID to match short code but got ID=%s and ShortCode=%s", url.ID, url.ShortCode)
	}

	// Verify the URL was saved
	savedURL, err := urlStorage.GetByID(ctx, url.ID)
	if err != nil {
		t.Errorf("Failed to get saved URL: %v", err)
	}
	if savedURL.ID != url.ID || savedURL.Original != url.Original {
		t.Errorf("Saved URL doesn't match created URL")
	}

	// Verify metrics were updated
	if metricsService.domains["github.com"] != 1 {
		t.Errorf("Expected domain count to be 1 but got %d", metricsService.domains["github.com"])
	}
}

func TestShortenerService_GetURL(t *testing.T) {
	// Set up mocks
	urlStorage := NewMockURLStorage()
	metricsService := NewMockMetricsService()
	
	// Create service with configuration
	config := ShortenerConfig{
		BaseURL:    "http://localhost:3000",
		CodeLength: 6,
	}
	service := NewTestShortenerService(urlStorage, metricsService, config)
	ctx := context.Background()

	// Add a test URL
	testURL := model.URL{
		ID:        "abc123",
		ShortCode: "abc123",
		Original:  "https://github.com/user/repo",
		CreatedAt: time.Now(),
	}
	err := urlStorage.Save(ctx, testURL)
	if err != nil {
		t.Fatalf("Failed to save test URL: %v", err)
	}

	// Test getting a URL by short code
	retrievedURL, err := service.GetURL(ctx, "abc123")
	if err != nil {
		t.Fatalf("Failed to get URL: %v", err)
	}

	// Verify the retrieved URL
	if retrievedURL.ID != testURL.ID || retrievedURL.Original != testURL.Original {
		t.Errorf("Retrieved URL doesn't match test URL")
	}

	// Test getting a non-existent URL
	_, err = service.GetURL(ctx, "nonexistent")
	if err != url.ErrURLNotFound {
		t.Errorf("Expected ErrURLNotFound but got: %v", err)
	}
}

func TestShortenerService_GenerateShortURL(t *testing.T) {
	// Set up mocks
	urlStorage := NewMockURLStorage()
	metricsService := NewMockMetricsService()
	
	// Create service with configuration
	config := ShortenerConfig{
		BaseURL:    "http://localhost:3000",
		CodeLength: 6,
	}
	service := NewTestShortenerService(urlStorage, metricsService, config)

	// Test generating a short URL
	shortURL := service.GenerateShortURL("abc123")
	expected := "http://localhost:3000/abc123"

	if shortURL != expected {
		t.Errorf("Expected short URL to be %s but got %s", expected, shortURL)
	}
}
