package searchengine

import (
	"sync"
)

var (
	model Model
)

// Model represents a data structure that manages an index, a store of documents, and maintains the order
// in which documents are inserted. It also includes a mutex for concurrent access.
type Model struct {
	// index is a map where the key is a string and the value is a slice of strings, typically used to store indexed data.
	index map[string][]string

	// store is a map where the key is a string and the value is a map of string to any type, representing a document store.
	store map[string]map[string]interface{}

	// order is a slice of strings representing the order in which documents were inserted or processed.
	order []string

	// mu is a mutex used to ensure thread-safety for concurrent access to the Model.
	mu *sync.Mutex
}

// MatchType defines the type of matching logic used for searches.
type MatchType string

// Sort defines the type of sorting used for search results.
type Sort string

// Constants for MatchType.
const (
	// MatchAND is used when a search should match documents where all conditions are true.
	MatchAND MatchType = "AND"

	// MatchOR is used when a search should match documents where any condition is true.
	MatchOR MatchType = "OR"

	// MatchPHRASE is used for phrase searches where the exact order of words matters.
	MatchPHRASE MatchType = "PHRASE"
)

// Constants for Sort.
const (
	// ASCENDING represents ascending sort order.
	ASCENDING Sort = "ascending"

	// DESCENDING represents descending sort order.
	DESCENDING Sort = "descending"
)

// DefaultOptions provides the default search options for Find operations.
var (
	// DefaultOptions defines the default values for searching, such as:
	// Start: 0 (starting index), Limit: 20 (number of results),
	// MatchType: MatchAND (using AND for matches), and Sort: DESCENDING (sorting results in descending order).
	DefaultOptions = FindOptions{
		Start:     0,
		Limit:     20,
		MatchType: MatchAND,
		Sort:      DESCENDING,
	}
)

// FindOptions defines the options used to customize search queries.
type FindOptions struct {
	// Start is the starting index for pagination (defaults to 0).
	Start int

	// Limit is the maximum number of results to return (defaults to 20).
	Limit int

	// MatchType defines how the search should match documents (AND, OR, PHRASE).
	MatchType MatchType

	// Sort defines the order in which search results are sorted (ascending or descending).
	Sort Sort
}
