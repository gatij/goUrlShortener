package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gatij/goUrlShortener/internal/api"
	"github.com/gatij/goUrlShortener/internal/service"
	"github.com/gatij/goUrlShortener/internal/storage/metrics"
	"github.com/gatij/goUrlShortener/internal/storage/url"
)

func setupTestRouter() http.Handler {
	// Initialize storage
	urlStore := url.NewMemoryStorage()
	metricsStore := metrics.NewMemoryStorage()

	// Initialize services
	metricsService := service.NewMetricsService(metricsStore)
	shortenerConfig := service.ShortenerConfig{
		BaseURL:    "http://localhost:3000",
		CodeLength: 6,
	}
	shortenerService := service.NewShortenerService(urlStore, metricsService, shortenerConfig)

	// Setup router
	return api.SetupRouter(shortenerService, metricsService)
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	
	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	
	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK' but got %s", w.Body.String())
	}
}

func TestCreateAndRetrieveURL(t *testing.T) {
	router := setupTestRouter()
	
	// Create a shortened URL
	createBody := map[string]string{
		"url": "https://github.com/golang/go",
	}
	jsonData, _ := json.Marshal(createBody)
	
	createReq, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)
	
	// Verify create response
	if createResp.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, createResp.Code)
	}
	
	// Parse response to get short code
	var createResult map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &createResult)
	
	shortCode, ok := createResult["short_code"].(string)
	if !ok || shortCode == "" {
		t.Fatalf("Failed to get short code from response: %v", createResult)
	}
	
	// Try to access the shortened URL
	redirectReq, _ := http.NewRequest("GET", "/"+shortCode, nil)
	redirectResp := httptest.NewRecorder()
	router.ServeHTTP(redirectResp, redirectReq)
	
	// Verify redirect
	if redirectResp.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status code %d but got %d", http.StatusMovedPermanently, redirectResp.Code)
	}
	
	location := redirectResp.Header().Get("Location")
	if location != "https://github.com/golang/go" {
		t.Errorf("Expected redirect to 'https://github.com/golang/go' but got %s", location)
	}
	
	// Check metrics
	metricsReq, _ := http.NewRequest("GET", "/api/v1/metrics/domains", nil)
	metricsResp := httptest.NewRecorder()
	router.ServeHTTP(metricsResp, metricsReq)
	
	// Verify metrics response
	if metricsResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, metricsResp.Code)
	}
	
	var metricsResult map[string]interface{}
	json.Unmarshal(metricsResp.Body.Bytes(), &metricsResult)
	
	topDomains, ok := metricsResult["top_domains"].([]interface{})
	if !ok || len(topDomains) == 0 {
		t.Errorf("Expected non-empty top_domains but got: %v", metricsResult)
	}
}

func TestRootEndpoint(t *testing.T) {
	router := setupTestRouter()
	
	// Create request
	req, _ := http.NewRequest("GET", "/", nil)
	
	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify service info is present
	if _, ok := response["service"]; !ok {
		t.Errorf("Expected 'service' field in response but got: %v", response)
	}
	
	if _, ok := response["endpoints"]; !ok {
		t.Errorf("Expected 'endpoints' field in response but got: %v", response)
	}
}
