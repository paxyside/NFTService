package utils

import (
	"regexp"
	"testing"
)

func TestGenerateUniqueHash(t *testing.T) {
	hash := GenerateUniqueHash()
	if len(hash) != 20 {
		t.Errorf("expected hash length to be 20, but got %d", len(hash))
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(hash) {
		t.Errorf("expected hash to only contain letters and digits, but got: %s", hash)
	}
}
