package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"royaka/internal/database"
	"royaka/internal/network"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := database.Connect(ctx, os.Getenv("DATABASE_URL")); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()
	if os.Getenv("LEGACY_SESSION_CUTOFF") == "" {
		log.Fatal("LEGACY_SESSION_CUTOFF is required (RFC3339 UTC timestamp)")
	}
	legacyCutoff, err := time.Parse(time.RFC3339, os.Getenv("LEGACY_SESSION_CUTOFF"))
	if err != nil {
		log.Fatalf("LEGACY_SESSION_CUTOFF must be RFC3339: %v", err)
	}
	network.ConfigureLegacySessionCutoff(legacyCutoff)
	if err := network.CleanupExpiredSessions(ctx); err != nil {
		log.Fatalf("Legacy session cleanup failed: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// WebSocket handler
	http.HandleFunc("/ws", network.HandleWebSocket)

	log.Println("Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
