package ideas

import (
	"testing"
)

func TestGetSecretNum(t *testing.T) {
	secret := getSecretNum()
	if len(secret) != numDigits {
		t.Errorf("Expected length %d, got %d", numDigits, len(secret))
	}

	for _, r := range secret {
		if r < '0' || r > '9' {
			t.Errorf("Non-digit in secret: %c", r)
		}
	}

	seen := make(map[rune]bool)
	for _, r := range secret {
		if seen[r] {
			t.Errorf("Duplicate digit: %c", r)
		}
		seen[r] = true
	}
}

func TestGetClues(t *testing.T) {
	tests := []struct {
		guess, secret, expected string
	}{
		{"123", "123", "You got it!"},
		{"123", "456", "Bagels"},
		{"123", "132", "Fermi Pico Pico"},
		{"167", "625", "Pico"},
		{"145", "625", "Fermi"},
		{"111", "625", "Bagels"},
		{"111", "123", "Fermi"},
		{"248", "843", "Fermi Pico"},
	}
	for _, tt := range tests {
		result := getClues(tt.guess, tt.secret)
		if result != tt.expected {
			t.Errorf("getClues(%q, %q) = %q; want %q", tt.guess, tt.secret, result, tt.expected)
		}
	}
}
