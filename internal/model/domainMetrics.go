package model

type DomainMetrics struct {
	Domain      string `json:"domain"`       // The domain name
	ShortenCount int    `json:"shorten_count"` // Number of times URLs were shortened for this domain
}