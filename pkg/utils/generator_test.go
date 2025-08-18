package utils

import (
	"testing"
	"regexp"
)

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		name         string
		length       int
		expectedLen  int
		expectError  bool
	}{
		{
			name:        "Default length",
			length:      0,
			expectedLen: DefaultShortCodeLength,
			expectError: false,
		},
		{
			name:        "Custom length",
			length:      8,
			expectedLen: 8,
			expectError: false,
		},
		{
			name:        "Very short length",
			length:      3,
			expectedLen: 3,
			expectError: false,
		},
	}

	// Regular expression to match the alphabet used in short codes
	validCharsRegex := regexp.MustCompile("^[" + DefaultAlphabet + "]+$")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := GenerateShortCode(tt.length)
			
			// Check for expected errors
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got nil")
				return
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}
			
			// Check length
			if len(code) != tt.expectedLen {
				t.Errorf("Expected length %d but got %d", tt.expectedLen, len(code))
			}
			
			// Check that only valid characters are used
			if !validCharsRegex.MatchString(code) {
				t.Errorf("Generated code %q contains invalid characters", code)
			}
			
			// Check uniqueness (generate multiple codes)
			if !tt.expectError {
				codes := make(map[string]bool)
				for i := 0; i < 100; i++ {
					newCode, _ := GenerateShortCode(tt.length)
					if codes[newCode] {
						t.Errorf("Generated duplicate code: %s", newCode)
					}
					codes[newCode] = true
				}
			}
		})
	}
}

func TestGenerateShortURL(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		shortCode string
		expected  string
	}{
		{
			name:      "Standard base URL",
			baseURL:   "https://short.io",
			shortCode: "abc123",
			expected:  "https://short.io/abc123",
		},
		{
			name:      "Base URL with trailing slash",
			baseURL:   "https://short.io/",
			shortCode: "abc123",
			expected:  "https://short.io/abc123",
		},
		{
			name:      "Localhost URL",
			baseURL:   "http://localhost:3000",
			shortCode: "xyz789",
			expected:  "http://localhost:3000/xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateShortURL(tt.baseURL, tt.shortCode)
			if result != tt.expected {
				t.Errorf("Expected %q but got %q", tt.expected, result)
			}
		})
	}
}
