package main

import (
	"server/config"
	"server/driver/csv"
	searchengine "server/driver/search-engine"
	"testing"
)

func TestInit(t *testing.T) {
	err := config.Init()
	if err != nil {
		t.Fatal()
	}
}

func BenchmarkParsing(b *testing.B) {
	searchengine.Init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := csv.ParseCSVFromDir("storage/upload")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	searchengine.Init()
	err := csv.ParseCSVFromDir("storage/upload")
	if err != nil {
		b.Fatal(err)
	}
	filter := map[string]interface{}{
		"Message": "error",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := searchengine.Find(filter)
		if err != nil {
			b.Fatal(err)
		}
	}
}
