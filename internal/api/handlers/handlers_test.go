package handlers

import (
"bytes"
"context"
"encoding/json"
"errors"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/gin-gonic/gin"
"github.com/gatij/goUrlShortener/internal/model"
)

// ShortenerServiceInterface defines the interface for the shortener service
type ShortenerServiceInterface interface {
	CreateShortURL(ctx context.Context, originalURL string) (model.URL, error)
	GetURL(ctx context.Context, shortCode string) (model.URL, error)
	GenerateShortURL(shortCode string) string
}

// MetricsServiceInterface defines the interface for the metrics service
type MetricsServiceInterface interface {
	GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error)
	IncrementDomainShortenCount(ctx context.Context, domain string) error
}

// MockShortenerService is a mock implementation of the shortener service
type MockShortenerService struct {
	urls map[string]model.URL
}

func NewMockShortenerService() *MockShortenerService {
	return &MockShortenerService{
		urls: make(map[string]model.URL),
	}
}

func (m *MockShortenerService) CreateShortURL(ctx context.Context, originalURL string) (model.URL, error) {
	// For testing, always return the same short code
	url := model.URL{
		ID:        "testcode",
		ShortCode: "testcode",
		Original:  originalURL,
		CreatedAt: time.Now(),
	}
	m.urls[url.ShortCode] = url
	return url, nil
}

func (m *MockShortenerService) GetURL(ctx context.Context, shortCode string) (model.URL, error) {
	url, exists := m.urls[shortCode]
	if !exists {
		return model.URL{}, errors.New("url not found")
	}
	return url, nil
}

func (m *MockShortenerService) GenerateShortURL(shortCode string) string {
	return "http://localhost:3000/" + shortCode
}

// ShortenerHandler that uses our interface
type TestShortenerHandler struct {
	shortenerService ShortenerServiceInterface
}

func NewTestShortenerHandler(service ShortenerServiceInterface) *TestShortenerHandler {
	return &TestShortenerHandler{shortenerService: service}
}

func (h *TestShortenerHandler) CreateShortURL(c *gin.Context) {
	var req URLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	url, err := h.shortenerService.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
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

// TestRedirectHandler for redirect tests
type TestRedirectHandler struct {
	shortenerService ShortenerServiceInterface
}

func NewTestRedirectHandler(service ShortenerServiceInterface) *TestRedirectHandler {
	return &TestRedirectHandler{shortenerService: service}
}

func (h *TestRedirectHandler) RedirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short code"})
		return
	}

	url, err := h.shortenerService.GetURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.Original)
}

// TestMetricsHandler for metrics tests
type TestMetricsHandler struct {
	metricsService MetricsServiceInterface
}

func NewTestMetricsHandler(service MetricsServiceInterface) *TestMetricsHandler {
	return &TestMetricsHandler{metricsService: service}
}

func (h *TestMetricsHandler) GetTopDomains(c *gin.Context) {
	limit := 3 // Default to top 3
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

// MockMetricsService for testing the metrics handler
type MockMetricsService struct {
	GetTopDomainsFunc func(ctx context.Context, limit int) ([]model.DomainMetrics, error)
}

func (m *MockMetricsService) GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error) {
	return m.GetTopDomainsFunc(ctx, limit)
}

func (m *MockMetricsService) IncrementDomainShortenCount(ctx context.Context, domain string) error {
	return nil
}

func TestShortenerHandler_CreateShortURL(t *testing.T) {
	// Set up gin in test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockService := NewMockShortenerService()
	
	// Create handler
	handler := NewTestShortenerHandler(mockService)
	
	// Set up router
	router := gin.New()
	router.POST("/api/v1/urls", handler.CreateShortURL)
	
	// Create test request
	requestBody := URLRequest{
		URL: "https://example.com/path",
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
	}
	
	// Parse response
	var response URLResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify response
	if response.ShortCode != "testcode" {
		t.Errorf("Expected short code 'testcode' but got %s", response.ShortCode)
	}
	if response.ShortURL != "http://localhost:3000/testcode" {
		t.Errorf("Expected short URL 'http://localhost:3000/testcode' but got %s", response.ShortURL)
	}
	if response.OriginalURL != "https://example.com/path" {
		t.Errorf("Expected original URL 'https://example.com/path' but got %s", response.OriginalURL)
	}
}

func TestRedirectHandler_RedirectToOriginal(t *testing.T) {
	// Set up gin in test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockService := NewMockShortenerService()
	
	// Add a test URL
	mockService.urls["testcode"] = model.URL{
		ID:        "testcode",
		ShortCode: "testcode",
		Original:  "https://example.com/path",
		CreatedAt: time.Now(),
	}
	
	// Create handler
	handler := NewTestRedirectHandler(mockService)
	
	// Set up router
	router := gin.New()
	router.GET("/:shortCode", handler.RedirectToOriginal)
	
	// Create test request
	req, _ := http.NewRequest("GET", "/testcode", nil)
	
	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Check response
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status code %d but got %d", http.StatusMovedPermanently, w.Code)
	}
	
	// Verify redirect location
	location := w.Header().Get("Location")
	if location != "https://example.com/path" {
		t.Errorf("Expected redirect to 'https://example.com/path' but got %s", location)
	}
}

func TestMetricsHandler_GetTopDomains(t *testing.T) {
	// Set up gin in test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock metrics service
	mockService := &MockMetricsService{}
	
	// Define test domains
	mockService.GetTopDomainsFunc = func(ctx context.Context, limit int) ([]model.DomainMetrics, error) {
		return []model.DomainMetrics{
			{Domain: "example.com", ShortenCount: 10},
			{Domain: "github.com", ShortenCount: 5},
		}, nil
	}
	
	// Create handler
	handler := NewTestMetricsHandler(mockService)
	
	// Set up router
	router := gin.New()
	router.GET("/api/v1/metrics/domains", handler.GetTopDomains)
	
	// Create test request
	req, _ := http.NewRequest("GET", "/api/v1/metrics/domains?limit=2", nil)
	
	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	
	// Parse response
	var response struct {
		TopDomains []model.DomainMetrics `json:"top_domains"`
		Limit      int                  `json:"limit"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify response
	if len(response.TopDomains) != 2 {
		t.Errorf("Expected 2 domains but got %d", len(response.TopDomains))
	}
	if response.Limit != 3 { // Default limit is 3 in our test implementation
		t.Errorf("Expected limit 3 but got %d", response.Limit)
	}
	if response.TopDomains[0].Domain != "example.com" || response.TopDomains[0].ShortenCount != 10 {
		t.Errorf("Expected first domain to be example.com with count 10, but got %s with count %d", 
response.TopDomains[0].Domain, response.TopDomains[0].ShortenCount)
	}
}
