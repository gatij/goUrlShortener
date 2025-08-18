package url

import (
    "context"

    "github.com/gatij/goUrlShortener/internal/model"
)

// Storage defines the interface for URL storage operations
type Storage interface {
    // Save stores a new shortened URL
    Save(ctx context.Context, url model.URL) error

    // GetByID retrieves a URL by its short ID
    GetByID(ctx context.Context, id string) (model.URL, error)

    // GetByShortCode retrieves a URL by its short code
    GetByShortCode(ctx context.Context, shortCode string) (model.URL, error)

    // Delete removes a URL from storage
    Delete(ctx context.Context, id string) error
}