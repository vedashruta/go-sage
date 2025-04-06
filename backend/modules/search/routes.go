package search

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/get", get)
	mux.HandleFunc("/getDoc", getDoc)
	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/search", search)
	mux.HandleFunc("/upload", upload)
}
