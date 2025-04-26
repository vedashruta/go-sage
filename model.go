package gosage

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

type Document struct{
	
}

type MatchType string

const (
	MatchAND    MatchType = "AND"
	MatchOR     MatchType = "OR"
	MatchPHRASE MatchType = "PHRASE"
)

var (
	DefaultOptions = FindOptions{
		Limit:     20,
		MatchType: MatchAND,
	}
)

type FindOptions struct {
	Limit     int
	MatchType MatchType
}
