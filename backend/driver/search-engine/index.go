package searchengine

import (
	"fmt"
	"sync"
	"time"
)

func Init() {
	model = Model{
		index: make(map[string][]string),
		store: make(map[string]map[string]interface{}),
		order: []string{},
		mu:    &sync.Mutex{},
	}
}

func Index(doc map[string]interface{}, docID string) (duration time.Duration, count int, err error) {
	start := time.Now()
	model.mu.Lock()
	defer model.mu.Unlock()

	// Overwrite the document
	model.store[docID] = doc

	// Remove the docID from model.order if it already exists
	for i, id := range model.order {
		if id == docID {
			model.order = append(model.order[:i], model.order[i+1:]...)
			break
		}
	}
	// Re-append to mark as latest
	model.order = append(model.order, docID)

	// Indexing
	for _, val := range doc {
		strVal := fmt.Sprintf("%v", val)
		tokens := stemTokens(removeStopWords(tokenize(normalize(strVal))))
		for _, term := range tokens {
			model.index[term] = append(model.index[term], docID)
			count++
		}
	}

	duration = time.Since(start)
	return
}
