package url

import (
    "context"
    "errors"
    "sync"

    "github.com/PuerkitoBio/purell"
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
    urls              map[string]model.URL  // Maps ID to URL object
    shortToURL        map[string]string     // Maps short code to ID
    normalizedToShort map[string]string     // Maps normalized URL to short code
    mu                sync.RWMutex          // Protects the maps from concurrent access
}

// NewMemoryStorage creates a new in-memory URL storage
func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        urls:              make(map[string]model.URL),
        shortToURL:        make(map[string]string),
        normalizedToShort: make(map[string]string),
    }
}

// normalizeURL standardizes a URL for consistent lookups using purell
func (s *MemoryStorage) normalizeURL(rawURL string) string {
    // Define normalization flags
    flags := purell.FlagsSafe | purell.FlagRemoveTrailingSlash | 
             purell.FlagRemoveDotSegments | purell.FlagRemoveDuplicateSlashes |
             purell.FlagSortQuery | purell.FlagRemoveEmptyQuerySeparator

    // Normalize the URL
    normalized, err := purell.NormalizeURLString(rawURL, flags)
    if err != nil {
        // If normalization fails, return the original URL
        return rawURL
    }
    
    return normalized
}

// Save stores a new shortened URL
func (s *MemoryStorage) Save(ctx context.Context, url model.URL) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Check if ID already exists
    if _, exists := s.urls[url.ID]; exists {
        return ErrURLExists
    }
    
    // Normalize the original URL
    normalizedURL := s.normalizeURL(url.Original)
    
    // Check if the normalized URL already exists
    if existingShort, exists := s.normalizedToShort[normalizedURL]; exists {
        // Return the existing URL object
        id := s.shortToURL[existingShort]
        existingURL := s.urls[id]
        
        // Update the provided URL object with existing values (for return by reference)
        url = existingURL
        
        return nil
    }
    
    // Store URL by ID
    s.urls[url.ID] = url
    
    // Store mapping from short code to ID
    s.shortToURL[url.ShortCode] = url.ID
    
    // Store mapping from normalized original URL to short code
    s.normalizedToShort[normalizedURL] = url.ShortCode
    
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
    id, exists := s.shortToURL[shortCode]
    if !exists {
        return model.URL{}, ErrURLNotFound
    }
    
    // Get URL by ID
    return s.urls[id], nil
}

// GetByOriginalURL retrieves a URL by its original URL
func (s *MemoryStorage) GetByOriginalURL(ctx context.Context, originalURL string) (model.URL, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // Normalize the URL for consistent lookup
    normalizedURL := s.normalizeURL(originalURL)
    
    // Direct lookup from normalized URL to short code - O(1)
    shortCode, exists := s.normalizedToShort[normalizedURL]
    if !exists {
        return model.URL{}, ErrURLNotFound
    }
    
    // Get URL using the short code
    id := s.shortToURL[shortCode]
    return s.urls[id], nil
}

// Delete removes a URL from storage
func (s *MemoryStorage) Delete(ctx context.Context, id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    url, exists := s.urls[id]
    if !exists {
        return ErrURLNotFound
    }
    
    // Get the short code and normalize the original URL before deleting
    shortCode := url.ShortCode
    normalizedURL := s.normalizeURL(url.Original)
    
    // Remove from all maps
    delete(s.urls, id)
    delete(s.shortToURL, shortCode)
    delete(s.normalizedToShort, normalizedURL)
    
    return nil
}