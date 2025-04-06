package parquet

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	searchengine "server/driver/search-engine"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

var totalDocsParsed atomic.Int64

// Parse processes all parquet files in a given directory.
// It walks through the directory, parses each parquet file concurrently using goroutines,
// and indexes the parsed documents into a search engine.
// The total number of documents parsed is printed at the end.
func Parse(dirPath string) (err error) {
	// A channel to pass file paths for processing
	fileChan := make(chan string, 10)

	// WaitGroup to synchronize the completion of all goroutines
	var wg sync.WaitGroup

	// Start 4 goroutines to concurrently process the parquet files
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Process files received from the fileChan
			for path := range fileChan {
				// Measure the time taken to parse the file
				start := time.Now()
				count, err := parseParquetFile(path)
				duration := time.Since(start)
				if err != nil {
					log.Printf("Error reading file: %s\n Error: %s", path, err.Error())
				}
				// Output the result of parsing the file
				fmt.Printf("Parsed file: %-40s | Time: %v | Docs: %d\n", filepath.Base(path), duration, count)
			}
		}()
	}

	// Walk through the directory and send file paths to the channel
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only send file paths, not directories
		if !info.IsDir() {
			fileChan <- path
		}
		return nil
	})

	// Close the file channel and wait for all goroutines to finish
	close(fileChan)
	wg.Wait()

	// Output the total number of documents parsed
	fmt.Printf("\nTotal documents loaded: %d\n", totalDocsParsed.Load())
	return
}

// parseParquetFile reads a parquet file, processes its rows, and indexes the documents into a search engine.
// It returns the number of documents successfully indexed and any error encountered during parsing.
func parseParquetFile(filePath string) (count int, err error) {
	// Open the parquet file
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return
	}
	defer fr.Close()

	// Create a parquet reader for the file
	pr, err := reader.NewParquetReader(fr, nil, 4)
	if err != nil {
		return
	}
	defer pr.ReadStop()

	// Get the schema of the parquet file
	schema := pr.SchemaHandler.ValueColumns
	total := int(pr.GetNumRows())
	batchSize := 1000
	docChan := make(chan map[string]interface{}, batchSize)

	// WaitGroup to synchronize document processing
	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine to process documents and index them into the search engine
	go func() {
		defer wg.Done()
		for doc := range docChan {
			// Generate a unique docID using MsgId and a timestamp
			msgID := fmt.Sprintf("%v", doc["MsgId"])
			if msgID == "" || msgID == "<nil>" {
				msgID = uuid.New().String()
			}
			docID := fmt.Sprintf("%s_%d", msgID, time.Now().UnixNano())

			// Index the document
			_, _, err := searchengine.Index(doc, docID)
			if err != nil {
				log.Printf("Error indexing doc: %v", err)
			} else {
				// Update the total number of documents parsed
				totalDocsParsed.Add(1)
				count++
			}
		}
	}()

	// Read the parquet file in batches and send rows to the docChan for indexing
	for i := 0; i < total; i += batchSize {
		if i+batchSize > total {
			batchSize = total - i
		}
		columns := make([][]interface{}, len(schema))
		// Read a batch of rows from the parquet file
		err = pr.Read(&columns)
		if err != nil {
			return
		}

		// Process each row and send it to the docChan
		for rowIndex := 0; rowIndex < len(columns[0]); rowIndex++ {
			row := make(map[string]interface{})
			for colIndex, colName := range schema {
				if rowIndex < len(columns[colIndex]) {
					// Clean the column name and assign the value to the row map
					colNameParts := strings.Split(colName, ".")
					cleanColName := colNameParts[len(colNameParts)-1]
					val := columns[colIndex][rowIndex]
					switch v := val.(type) {
					case string:
						row[cleanColName] = v
					case int32, int64, float32, float64, bool:
						row[cleanColName] = fmt.Sprintf("%v", v)
					default:
						// Log unsupported types
						log.Printf("Unsupported type for field %s: %T, skipping", cleanColName, v)
					}
				}
			}
			// Send the row to the docChan for indexing
			docChan <- row
		}
	}
	// Close the docChan and wait for the goroutine to finish
	close(docChan)
	wg.Wait()
	return
}
