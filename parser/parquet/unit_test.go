package parquet

import (
	"testing"

	gosage "github.com/vedashruta/go-sage.git"
)

func BenchmarkParsing(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ParseParquetFromDir("../../storage/upload")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	err := ParseParquetFromDir("../../storage/upload")
	if err != nil {
		b.Fatal(err)
	}
	filter := map[string]interface{}{
		"Message": "error",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gosage.Find(filter)
		if err != nil {
			b.Fatal(err)
		}
	}
}
