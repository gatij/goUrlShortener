package model

import "time"

// URL represents a shortened URL entry in the system
type URL struct {
	ID        string    `json:"id"`        // Unique identifier for the URL
	ShortCode string    `json:"short_code"` // Shortened code for the URL
	Original  string    `json:"original"`  // Original long URL
	CreatedAt time.Time `json:"created_at"` // Timestamp when the URL was created
}