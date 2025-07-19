package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"olympus/routes"
	"olympus/telemetry"

	"github.com/gorilla/mux"
)

func main() {
	// Inisialisasi OpenTelemetry
	shutdown := telemetry.InitTracer()
	defer shutdown(context.Background())

	r := mux.NewRouter()
	routes.ApiV1Routes(r)

	// Setup HTTP Server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r, // <-- Pasang router
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Channel untuk menunggu sinyal interupsi
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		// Graceful shutdown dengan timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fmt.Println("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("HTTP server Shutdown: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	fmt.Println("Server running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("HTTP server ListenAndServe: %v\n", err)
	}

	<-idleConnsClosed
	fmt.Println("Server gracefully stopped")
}
