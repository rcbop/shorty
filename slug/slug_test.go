package slug

import (
	"testing"
)

func TestGenerateRandomSlug(t *testing.T) {
	slug := GenerateRandomSlug()

	if len(slug) != 6 {
		t.Errorf("Expected slug length of 6, but got %d", len(slug))
	}

	validCharset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range slug {
		if !isValidChar(char, validCharset) {
			t.Errorf("Invalid character in slug: %c", char)
		}
	}
}

func isValidChar(char rune, charset string) bool {
	for _, validChar := range charset {
		if validChar == char {
			return true
		}
	}
	return false
}
