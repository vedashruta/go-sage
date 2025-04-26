package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	gosage "github.com/vedashruta/go-sage.git"
)

var totalDocsParsed atomic.Int64

func ParseCSVFromDir(dirPath string) error {
	startAll := time.Now()

	fileChan := make(chan string, 10)
	var wg sync.WaitGroup

	// Launch worker goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				start := time.Now()

				file, err := os.Open(path)
				if err != nil {
					log.Printf("Failed to open CSV: %s | %v", path, err)
					continue
				}
				count, err := ParseFile(file)
				file.Close()

				if err != nil {
					log.Printf("Failed to parse CSV: %s | %v", path, err)
					continue
				}

				fmt.Printf("Parsed: %-40s | Time: %-10v | Docs: %d\n", filepath.Base(path), time.Since(start), count)
				totalDocsParsed.Add(int64(count))
			}
		}()
	}

	// Walk and send .csv files
	go func() {
		defer close(fileChan)
		filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".csv") {
				fileChan <- path
			}
			return nil
		})
	}()

	wg.Wait()
	totalTime := time.Since(startAll)
	fmt.Printf("\nTotal Documents Parsed: %d\n", totalDocsParsed.Load())
	fmt.Printf("Total Time Taken       : %v\n", totalTime)
	return nil
}

func ParseFile(r io.Reader) (count int, err error) {
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return
	}
	count = 0
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

		if err := indexDocument(row, docID); err != nil {
			continue
		}
		count++
	}
	return
}

func indexDocument(row map[string]interface{}, docID string) (err error) {
	err = gosage.Index(row, docID)
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
