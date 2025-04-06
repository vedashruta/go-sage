package router

import (
	"net/http"
	"server/modules/search"
)

// NewRouter initializes a new HTTP request multiplexer (router).
// It registers all search-related routes and a default handler for the root path.
// Returns the configured http.Handler to be used by the server.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Register search engine-related API routes.
	search.RegisterRoutes(mux)

	// Handle root path with an empty response or basic health check if needed.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Optional: Add a simple health check response here
	})

	return mux
}
