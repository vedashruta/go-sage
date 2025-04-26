package services

import (
	"fmt"
	"path/filepath"
	"time"
)

// Print logs a formatted summary line showing the base filename,
// time taken to parse it, and the number of documents parsed.
func Print(filePath string, duration time.Duration, docCount int) {
	fmt.Printf("Parsed: %-40s | Time: %-10v | Docs: %-4d\n", filepath.Base(filePath), duration, docCount)
}
