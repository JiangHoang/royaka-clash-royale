package main

import (
	"path/filepath"
	"testing"
)

func TestLegacyFixturesCanBeDecoded(t *testing.T) {
	dataDir := filepath.Join("..", "..", "assets", "data")
	var users []legacyUser
	var sessions []legacySession
	var troops []troop
	var towers []tower
	tests := []struct {
		name   string
		target any
		count  int
	}{
		{"users.json", &users, 13},
		{"sessions.json", &sessions, 13},
		{"troops.json", &troops, 60},
		{"towers.json", &towers, 2},
	}
	for _, test := range tests {
		if err := readJSON(filepath.Join(dataDir, test.name), test.target); err != nil {
			t.Fatalf("decode %s: %v", test.name, err)
		}
	}
	if len(users) != 13 || len(sessions) != 13 || len(troops) != 60 || len(towers) != 2 {
		t.Fatalf("unexpected fixture counts: users=%d sessions=%d troops=%d towers=%d",
			len(users), len(sessions), len(troops), len(towers))
	}
	if users[0].Password == "" || troops[0].Name == "" || towers[0].Type == "" {
		t.Fatal("important fixture fields were not decoded")
	}
}
