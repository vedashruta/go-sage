package main

import (
	"fmt"
	"log"
	"net/http"
	"server/config"
	"server/driver/csv"
	searchengine "server/driver/search-engine"
	"server/router"
	"time"
)

// main is the entry point of the Go Sage backend application.
// It initializes the configuration, sets up the search engine, parses CSV files from a directory,
// configures CORS, and starts the HTTP server.
func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the in-memory search engine model and store.
	searchengine.Init()

	dirPath := "storage/upload"

	// Parse all CSV files in the specified directory to populate the search engine.
	err = csv.ParseCSVFromDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Wrap the router with CORS middleware.
	handlerWithCORS := withCORS(router.NewRouter())

	// Set up and start the HTTP server using the configured host, port, and timeouts.
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port),
		Handler:      handlerWithCORS,
		ReadTimeout:  time.Duration(config.Config.ReadTimeout),
		WriteTimeout: time.Duration(config.Config.WriteTimeout),
	}

	fmt.Printf("Server running on %s:%d\n", config.Config.Host, config.Config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

// withCORS wraps an HTTP handler with CORS headers.
// This is primarily used for development to allow requests from any origin.
// In production, it's recommended to restrict allowed origins.
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
