package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"server/driver/helper"
	searchengine "server/driver/search-engine"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// ParseCSVFromDir walks through the given directory path, identifies all `.csv` files,
// and concurrently parses them using a worker pool of 10 goroutines.
// It prints the parsing duration and document count per file and finally summarizes total stats.
func ParseCSVFromDir(dirPath string) (err error) {
	startAll := time.Now()

	fileChan := make(chan string, 10) // Buffered channel to hold file paths
	var wg sync.WaitGroup

	var totalDocsParsed atomic.Int64 // Thread-safe counter

	// Launch worker goroutines to process CSV files concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				file, err := os.Open(path)
				if err != nil {
					log.Printf("Failed to open CSV: %s | %v", path, err)
					continue
				}
				docCount, duration, err := ParseFile(file)
				_ = file.Close() // Explicit close to avoid defer accumulation

				if err != nil {
					log.Printf("Failed to parse CSV: %s | %v", path, err)
					continue
				}
				helper.Print(path, duration, docCount)
				totalDocsParsed.Add(int64(docCount))
			}
		}()
	}

	// Walk the directory tree and send all `.csv` file paths into the fileChan
	go func() {
		defer close(fileChan)
		_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".csv") {
				fileChan <- path
			}
			return nil
		})
	}()

	wg.Wait() // Wait for all goroutines to finish
	totalTime := time.Since(startAll)
	fmt.Printf("\nTotal Documents Parsed : %d\n", totalDocsParsed.Load())
	fmt.Printf("Total Time Taken       : %v\n", totalTime)
	return
}

// ParseFile reads CSV content from the provided io.Reader, parses each row,
// converts it into a map[string]interface{}, and indexes it into the search engine.
// It returns the number of documents parsed and the total duration taken.
func ParseFile(r io.Reader) (docCount int, duration time.Duration, err error) {
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return
	}
	start := time.Now()
	docCount = 0

	for {
		record, err := reader.Read()
		if err != nil {
			break // Exit loop on EOF or read error
		}

		row := make(map[string]interface{})
		for i, h := range headers {
			key := cleanString(h)
			val := ""
			if i < len(record) {
				val = cleanString(record[i])
			}
			row[key] = val
		}

		// Use MsgId if available; otherwise, generate a new UUID
		msgID := fmt.Sprintf("%v", row["MsgId"])
		if msgID == "" {
			msgID = uuid.New().String()
		}
		docID := msgID

		_, _, err = searchengine.Index(row, docID)
		if err != nil {
			continue // Skip failed index attempts
		}
		docCount++
	}
	duration = time.Since(start)
	return
}

// cleanString sanitizes input strings by removing escape characters, extra quotes,
// trimming whitespace, and removing leading colons if present.
func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\\\"", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, ":")
	return s
}
