package network

import (
	"strings"
	"testing"
)

func TestUsernameValidation(t *testing.T) {
	for _, username := range []string{"a", "jiang@example.com", "Nguyễn Giang", "🎮"} {
		if !validUsername(username) {
			t.Errorf("expected %q to be valid", username)
		}
	}
	for _, username := range []string{"", "   ", strings.Repeat("a", 33)} {
		if validUsername(username) {
			t.Errorf("expected %q to be invalid", username)
		}
	}
}
