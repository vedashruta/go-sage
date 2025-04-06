package helper

import (
	"fmt"
	"path/filepath"
	"time"
)

func Print(filePath string, duration time.Duration, docCount int) {
	fmt.Printf("Parsed: %-40s | Time: %-10v | Docs: %-4d\n", filepath.Base(filePath), duration, docCount)
}
