package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"royaka/internal/database"
	"royaka/internal/network"
	"syscall"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

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

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/healthz", healthHandler)

	// WebSocket handler
	mux.HandleFunc("/ws", network.HandleWebSocket)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	log.Println("Server running at http://localhost:" + port)
	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	case <-shutdownSignal.Done():
		log.Println("Shutdown signal received")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			_ = server.Close()
		}
	}
}
