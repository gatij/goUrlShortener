package url

import (
    "context"
    "errors"
    "sync"

    "github.com/gatij/goUrlShortener/internal/model"
)

var (
    // ErrURLNotFound is returned when a URL is not found in storage
    ErrURLNotFound = errors.New("url not found")
    
    // ErrURLExists is returned when attempting to save a URL that already exists
    ErrURLExists = errors.New("url with this ID already exists")
)

// MemoryStorage implements the Storage interface with in-memory data structures
type MemoryStorage struct {
    urls      map[string]model.URL  // Maps ID to URL object
    shortCodes map[string]string    // Maps short code to ID for quick lookups
    mu        sync.RWMutex          // Protects the maps from concurrent access
}

// NewMemoryStorage creates a new in-memory URL storage
func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        urls:      make(map[string]model.URL),
        shortCodes: make(map[string]string),
    }
}

// Save stores a new shortened URL
func (s *MemoryStorage) Save(ctx context.Context, url model.URL) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Check if ID already exists
    if _, exists := s.urls[url.ID]; exists {
        return ErrURLExists
    }
    
    // Store URL by ID
    s.urls[url.ID] = url
    
    // Store mapping from short code to ID for quick lookups
    // Note: In this implementation, we're assuming the short code is the same as ID
    // If they're different in your model, adjust accordingly
    s.shortCodes[url.ID] = url.ID
    
    return nil
}

// GetByID retrieves a URL by its ID
func (s *MemoryStorage) GetByID(ctx context.Context, id string) (model.URL, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    url, exists := s.urls[id]
    if !exists {
        return model.URL{}, ErrURLNotFound
    }
    
    return url, nil
}

// GetByShortCode retrieves a URL by its short code
func (s *MemoryStorage) GetByShortCode(ctx context.Context, shortCode string) (model.URL, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // Look up ID from short code
    id, exists := s.shortCodes[shortCode]
    if !exists {
        return model.URL{}, ErrURLNotFound
    }
    
    // Get URL by ID
    return s.GetByID(ctx, id)
}

// Delete removes a URL from storage
func (s *MemoryStorage) Delete(ctx context.Context, id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    url, exists := s.urls[id]
    if !exists {
        return ErrURLNotFound
    }
    
    // Remove from both maps
    delete(s.urls, id)
    delete(s.shortCodes, url.ID) // Assuming ID is the short code
    
    return nil
}