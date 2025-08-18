package url

import (
	"context"
	"testing"
	"time"

	"github.com/gatij/goUrlShortener/internal/model"
)

func TestMemoryStorage_Save(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create test URL
	url := model.URL{
		ID:        "abc123",
		ShortCode: "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}

	// Test saving a new URL
	err := storage.Save(ctx, url)
	if err != nil {
		t.Errorf("Failed to save URL: %v", err)
	}

	// Test saving a duplicate URL
	err = storage.Save(ctx, url)
	if err != ErrURLExists {
		t.Errorf("Expected ErrURLExists but got: %v", err)
	}
}

func TestMemoryStorage_GetByID(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create and save test URL
	expectedURL := model.URL{
		ID:        "abc123",
		ShortCode: "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}

	err := storage.Save(ctx, expectedURL)
	if err != nil {
		t.Fatalf("Failed to save URL: %v", err)
	}

	// Test getting existing URL
	url, err := storage.GetByID(ctx, "abc123")
	if err != nil {
		t.Errorf("Failed to get URL by ID: %v", err)
	}

	if url.ID != expectedURL.ID || url.Original != expectedURL.Original {
		t.Errorf("Expected URL %+v but got %+v", expectedURL, url)
	}

	// Test getting non-existent URL
	_, err = storage.GetByID(ctx, "nonexistent")
	if err != ErrURLNotFound {
		t.Errorf("Expected ErrURLNotFound but got: %v", err)
	}
}

func TestMemoryStorage_GetByShortCode(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create and save test URL
	expectedURL := model.URL{
		ID:        "abc123",
		ShortCode: "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}

	err := storage.Save(ctx, expectedURL)
	if err != nil {
		t.Fatalf("Failed to save URL: %v", err)
	}

	// Test getting existing URL by short code
	url, err := storage.GetByShortCode(ctx, "abc123")
	if err != nil {
		t.Errorf("Failed to get URL by short code: %v", err)
	}

	if url.ID != expectedURL.ID || url.Original != expectedURL.Original {
		t.Errorf("Expected URL %+v but got %+v", expectedURL, url)
	}

	// Test getting non-existent URL by short code
	_, err = storage.GetByShortCode(ctx, "nonexistent")
	if err != ErrURLNotFound {
		t.Errorf("Expected ErrURLNotFound but got: %v", err)
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()

	// Create and save test URL
	url := model.URL{
		ID:        "abc123",
		ShortCode: "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}

	err := storage.Save(ctx, url)
	if err != nil {
		t.Fatalf("Failed to save URL: %v", err)
	}

	// Test deleting existing URL
	err = storage.Delete(ctx, "abc123")
	if err != nil {
		t.Errorf("Failed to delete URL: %v", err)
	}

	// Verify URL was deleted
	_, err = storage.GetByID(ctx, "abc123")
	if err != ErrURLNotFound {
		t.Errorf("URL was not deleted properly")
	}

	// Test deleting non-existent URL
	err = storage.Delete(ctx, "nonexistent")
	if err != ErrURLNotFound {
		t.Errorf("Expected ErrURLNotFound but got: %v", err)
	}
}
