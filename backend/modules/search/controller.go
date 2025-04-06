package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/driver/csv"
	"server/driver/helper"
	searchengine "server/driver/search-engine"
	"strconv"

	"github.com/google/uuid"
)

// get handles GET /get requests.
// It fetches a limited number of documents from the search engine.
func get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var limit int
	sizeStr := r.URL.Query().Get("size")
	if sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			limit = size
		}
	}
	if limit < 0 {
		limit = 20
	}
	docs := searchengine.Get(&searchengine.FindOptions{Limit: limit})
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(docs)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// getDoc handles GET /getDoc requests with a `query` parameter.
// It fetches matching documents from the search engine.
func getDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var limit int
	queryStr := r.URL.Query().Get("query")
	if queryStr == "" {
		w.Header().Set("Content-Type", "application/json")
		info := map[string]string{
			"error": "query key not found",
		}
		err := json.NewEncoder(w).Encode(info)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
	sizeStr := r.URL.Query().Get("size")
	if sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			limit = size
		}
	}
	if limit < 0 {
		limit = 20
	}
	docs := searchengine.GetDoc(queryStr, &searchengine.FindOptions{Limit: limit})
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(docs)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// stats handles GET /stats requests.
// It returns total indexed document count.
func stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	total := searchengine.GetStats()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]int{"total": total})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// search handles POST /search requests.
// It processes complex search queries with filters and pagination.
func search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Query     map[string]interface{} `json:"query"`
		Limit     int                    `json:"limit"`
		Start     int                    `json:"start"`
		Sort      string                 `json:"sort"`      // "ascending" or "descending"
		MatchType string                 `json:"matchType"` // Currently unused unless extended in Find
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var sortOrder searchengine.Sort
	switch req.Sort {
	case "ascending":
		sortOrder = searchengine.ASCENDING
	case "descending":
		sortOrder = searchengine.DESCENDING
	default:
		sortOrder = searchengine.DESCENDING
	}
	options := &searchengine.FindOptions{
		Limit: req.Limit,
		Start: req.Start,
		Sort:  sortOrder,
	}
	results, err := searchengine.Find(req.Query, options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// upload handles file uploads to the server.
// It supports .csv and .parquet files, saves them locally, and parses them.
func upload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate extension
	ext := filepath.Ext(header.Filename)
	if ext != ".csv" && ext != ".parquet" {
		http.Error(w, "Only .csv and .parquet files are allowed", http.StatusBadRequest)
		return
	}

	// Create local copy
	uuid := uuid.NewString()
	savePath := fmt.Sprintf("%[1]s/%[2]s_%[3]s", config.Config.Storage, uuid, header.Filename)
	outFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Create a TeeReader so we can write to disk and pass to ParseFile
	var parseBuffer bytes.Buffer
	tee := io.TeeReader(file, &parseBuffer)

	// Write to disk
	if _, err := io.Copy(outFile, tee); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Parse in-memory
	count, duration, err := csv.ParseFile(&parseBuffer)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse failed: %v", err), http.StatusInternalServerError)
		return
	}
	helper.Print(savePath, duration, count)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "File saved and parsed successfully",
		"documentsParsed": count,
		"duration":        duration.String(),
	})
}
