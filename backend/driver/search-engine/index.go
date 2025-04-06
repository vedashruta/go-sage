package searchengine

import (
	"fmt"
	"sync"
	"time"
)

// Init initializes the global in-memory model that holds the index, document store, and ordering.
func Init() {
	model = Model{
		index: make(map[string][]string),
		store: make(map[string]map[string]interface{}),
		order: []string{},
		mu:    &sync.Mutex{},
	}
}

// Index indexes a given document into the in-memory model using the specified docID.
// It normalizes, tokenizes, removes stop words, and stems the document fields,
// then builds an inverted index for fast lookup.
// Returns the duration of indexing, number of tokens indexed, and any error encountered.
func Index(doc map[string]interface{}, docID string) (duration time.Duration, count int, err error) {
	start := time.Now()
	model.mu.Lock()
	defer model.mu.Unlock()

	// Store the document (overwrite if it already exists)
	model.store[docID] = doc

	// Remove the docID from ordering if it already exists to avoid duplicates
	for i, id := range model.order {
		if id == docID {
			model.order = append(model.order[:i], model.order[i+1:]...)
			break
		}
	}

	// Append the docID to the end of the order list
	model.order = append(model.order, docID)

	// Build the inverted index for the document
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
