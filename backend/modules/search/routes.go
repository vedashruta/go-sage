package search

import "net/http"

// RegisterRoutes registers the HTTP handlers for various search engine endpoints.
//
// The following routes are registered:
// - /get     : Returns a list of recent documents.
// - /getDoc  : Retrieves a specific document by its ID.
// - /stats   : Returns statistics related to the indexed documents.
// - /search  : Performs a search based on query filters.
// - /upload  : Handles document uploads (e.g., CSV or Parquet files).
//
// It attaches these handlers to the provided ServeMux instance.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/get", get)
	mux.HandleFunc("/getDoc", getDoc)
	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/search", search)
	mux.HandleFunc("/upload", upload)
}
