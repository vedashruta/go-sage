package searchengine

import (
	"sync"
)

var (
	model Model
)

type Model struct {
	index map[string][]string
	store map[string]map[string]interface{}
	order []string
	mu    *sync.Mutex
}

type MatchType string
type Sort string

const (
	MatchAND    MatchType = "AND"
	MatchOR     MatchType = "OR"
	MatchPHRASE MatchType = "PHRASE"
	ASCENDING   Sort      = "ascending"
	DESCENDING  Sort      = "descending"
)

var (
	DefaultOptions = FindOptions{
		Start:     0,
		Limit:     20,
		MatchType: MatchAND,
		Sort:      DESCENDING,
	}
)

type FindOptions struct {
	Start     int
	Limit     int
	MatchType MatchType
	Sort      Sort
}
