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

func ParseCSVFromDir(dirPath string) (err error) {
	startAll := time.Now()

	fileChan := make(chan string, 10)
	var wg sync.WaitGroup

	var totalDocsParsed atomic.Int64

	// Launch worker goroutines
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
				_ = file.Close() // Explicit close to avoid defer leak

				if err != nil {
					log.Printf("Failed to parse CSV: %s | %v", path, err)
					continue
				}
				helper.Print(path, duration, docCount)
				totalDocsParsed.Add(int64(docCount))
			}
		}()
	}

	// Walk and send .csv files
	go func() {
		defer close(fileChan)
		_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".csv") {
				fileChan <- path
			}
			return nil
		})
	}()
	wg.Wait()
	totalTime := time.Since(startAll)
	fmt.Printf("\nTotal Documents Parsed : %d\n", totalDocsParsed.Load())
	fmt.Printf("Total Time Taken       : %v\n", totalTime)
	return
}

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
			break
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

		msgID := fmt.Sprintf("%v", row["MsgId"])
		if msgID == "" {
			msgID = uuid.New().String()
		}
		docID := msgID

		_, _, err = searchengine.Index(row, docID)
		if err != nil {
			continue
		}
		docCount++
	}
	duration = time.Since(start)
	return
}

func cleanString(s string) string {
	// Replace escaped quotes, extra colons, or malformed keys
	s = strings.ReplaceAll(s, "\\\"", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, ":")
	return s
}
