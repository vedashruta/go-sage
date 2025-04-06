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

func Parse(dirPath string) (err error) {
	fileChan := make(chan string, 10)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				start := time.Now()
				count, err := parseParquetFile(path)
				duration := time.Since(start)
				if err != nil {
					log.Printf("Error reading file: %s\n Error: %s", path, err.Error())
				}
				fmt.Printf("Parsed file: %-40s | Time: %v | Docs: %d\n", filepath.Base(path), duration, count)
			}
		}()
	}
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileChan <- path
		}
		return nil
	})
	close(fileChan)
	wg.Wait()
	fmt.Printf("\nTotal documents loaded: %d\n", totalDocsParsed.Load())
	return
}

func parseParquetFile(filePath string) (count int, err error) {
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, nil, 4)
	if err != nil {
		return
	}
	defer pr.ReadStop()

	schema := pr.SchemaHandler.ValueColumns
	total := int(pr.GetNumRows())
	batchSize := 1000
	docChan := make(chan map[string]interface{}, batchSize)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for doc := range docChan {
			// Generate docID using MsgId + random suffix or index fallback
			msgID := fmt.Sprintf("%v", doc["MsgId"])
			if msgID == "" || msgID == "<nil>" {
				msgID = uuid.New().String()
			}
			docID := fmt.Sprintf("%s_%d", msgID, time.Now().UnixNano())
			_, _, err := searchengine.Index(doc, docID)
			if err != nil {
				log.Printf("Error indexing doc: %v", err)
			} else {
				totalDocsParsed.Add(1)
				count++
			}
		}
	}()

	for i := 0; i < total; i += batchSize {
		if i+batchSize > total {
			batchSize = total - i
		}
		columns := make([][]interface{}, len(schema))
		err = pr.Read(&columns)
		if err != nil {
			return
		}
		for rowIndex := 0; rowIndex < len(columns[0]); rowIndex++ {
			row := make(map[string]interface{})
			for colIndex, colName := range schema {
				if rowIndex < len(columns[colIndex]) {
					colNameParts := strings.Split(colName, ".")
					cleanColName := colNameParts[len(colNameParts)-1]
					val := columns[colIndex][rowIndex]
					switch v := val.(type) {
					case string:
						row[cleanColName] = v
					case int32, int64, float32, float64, bool:
						row[cleanColName] = fmt.Sprintf("%v", v)
					default:
						log.Printf("Unsupported type for field %s: %T, skipping", cleanColName, v)
					}
				}
			}
			docChan <- row
		}
	}
	close(docChan)
	wg.Wait()
	return
}
