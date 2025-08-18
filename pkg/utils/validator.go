package utils

import (
    "errors"
    "net/url"
    "strings"

    "github.com/asaskevich/govalidator"
)

var (
    // ErrInvalidURL is returned when the URL is invalid
    ErrInvalidURL = errors.New("invalid URL format")
    
    // ErrBlockedDomain is returned when the URL domain is blocked
    ErrBlockedDomain = errors.New("domain is blocked")
    
    // ErrInfiniteLoopDetected is returned when attempting to shorten our own URL
    ErrInfiniteLoopDetected = errors.New("cannot shorten URLs from this service (infinite loop)")
    
    // ErrHTTPSRequired is returned when the URL doesn't use HTTPS
    ErrHTTPSRequired = errors.New("URL must use HTTPS protocol")
)

// List of domains that should not be shortened (to prevent redirect loops)
var selfDomains = []string{
    "localhost",
    "127.0.0.1",
    // Add your actual domain name(s) here
}

// URLInfo contains processed URL information
type URLInfo struct {
    OriginalURL   string // The original input URL
    NormalizedURL string // URL with enforced HTTPS if needed
    Domain        string // Extracted domain
}

// ValidateURL checks if a URL is valid and meets requirements
func ValidateURL(rawURL string) (*url.URL, error) {
    // Basic URL validation
    if !govalidator.IsURL(rawURL) {
        return nil, ErrInvalidURL
    }
    
    // Parse the URL
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return nil, ErrInvalidURL
    }
    
    // Check for infinite loop
    for _, domain := range selfDomains {
        if strings.Contains(parsedURL.Host, domain) {
            return nil, ErrInfiniteLoopDetected
        }
    }
    
    // Check for blocked domains
    blockedDomains := []string{"example.com", "malicious.com"}
    for _, blocked := range blockedDomains {
        if strings.Contains(parsedURL.Host, blocked) {
            return nil, ErrBlockedDomain
        }
    }
    
    // Validate scheme
    if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
        return nil, ErrInvalidURL
    }
    
    return parsedURL, nil
}

// ExtractDomain extracts the domain from a parsed URL
func ExtractDomain(parsedURL *url.URL) string {
    return parsedURL.Host
}

// EnforceHTTPS ensures the URL uses HTTPS, converting if necessary
func EnforceHTTPS(parsedURL *url.URL, enforceHTTPS bool) (string, error) {
    if parsedURL.Scheme != "https" {
        if parsedURL.Scheme == "http" && enforceHTTPS {
            // Convert to HTTPS if requested
            parsedURL.Scheme = "https"
            return parsedURL.String(), nil
        } else {
            return "", ErrHTTPSRequired
        }
    }
    
    // Already HTTPS
    return parsedURL.String(), nil
}

// ProcessURL orchestrates validation, domain extraction, and HTTPS enforcement
func ProcessURL(rawURL string, enforceHTTPS bool) (URLInfo, error) {
    info := URLInfo{
        OriginalURL: rawURL,
    }
    
    // Step 1: Validate URL
    parsedURL, err := ValidateURL(rawURL)
    if err != nil {
        return info, err
    }

	// Step 2: Enforce HTTPS if needed
    normalizedURL, err := EnforceHTTPS(parsedURL, enforceHTTPS)
    if err != nil {
        return info, err
    }
    info.NormalizedURL = normalizedURL
    
    // Step 3: Extract domain
    info.Domain = ExtractDomain(parsedURL)
    
    return info, nil
}

// IsValidShortCode checks if a short code meets our requirements
func IsValidShortCode(code string) bool {
    return govalidator.IsAlphanumeric(code) && len(code) >= 4 && len(code) <= 10
}