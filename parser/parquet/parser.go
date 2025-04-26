package parquet

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	gosage "github.com/vedashruta/go-sage.git"
	"github.com/vedashruta/go-sage.git/services"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

var totalDocsParsed atomic.Int64

// Parse processes all parquet files in a given directory.
// It walks through the directory, parses each parquet file concurrently using goroutines,
// and indexes the parsed documents into a search engine.
// The total number of documents parsed is printed at the end.
func ParseParquetFromDir(dirPath string) (err error) {
	startAll := time.Now()
	fileChan := make(chan string, 10) // Buffered channel to hold file paths
	var wg sync.WaitGroup
	// Launch worker goroutines to process CSV files concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				docCount, duration, err := ParseFile(path)
				if err != nil {
					log.Printf("Failed to parse CSV: %s | %v", path, err)
					continue
				}
				services.Print(path, duration, docCount)
				totalDocsParsed.Add(int64(docCount))
			}
		}()
	}

	// Walk the directory tree and send all `.csv` file paths into the fileChan
	go func() {
		defer close(fileChan)
		_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".parquet") {
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

// parseParquetFile reads a parquet file, processes its rows, and indexes the documents into a search engine.
// It returns the number of documents successfully indexed and any error encountered during parsing.
func ParseFile(filePath string) (docCount int, duration time.Duration, err error) {
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return
	}
	defer fr.Close()
	start := time.Now()
	docCount = 0
	pr, err := reader.NewParquetReader(fr, new(LogEntry), 10)
	if err != nil {
		return
	}
	defer pr.ReadStop()

	num := int(pr.GetNumRows())
	for i := 0; i < num; i += 10 {
		count := 10
		if num-i < 10 {
			count = num - i
		}
		var res []interface{}
		res, err = pr.ReadByNumber(count)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}
		for _, row := range res {
			if entry, ok := row.(LogEntry); ok {
				res := convertLogsToMap(entry)
				docID := res["MsgId"]
				err = gosage.Index(res, docID.(string))
				if err != nil {
					continue // Skip failed index attempts
				}
				docCount++
			} else {
				err = fmt.Errorf("unexpected row type: %T", row)
				return
			}
		}
	}
	duration = time.Since(start)
	return
}

func convertLogsToMap(input LogEntry) (res map[string]interface{}) {
	res = make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("json")
		fieldValue := val.Field(i).Interface()
		res[fieldName] = fieldValue
	}
	return
}
