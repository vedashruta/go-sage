package searchengine

// stopWords is a set of common English words that are excluded from indexing and search
// to reduce noise and improve relevance.
var stopWords = map[string]struct{}{
	"the": {}, "is": {}, "and": {}, "a": {}, "an": {}, "in": {}, "of": {}, "on": {}, "at": {}, "to": {},
}

// removeStopWords filters out common stop words from the list of tokens.
// It returns a new slice containing only the tokens that are not in the stopWords set.
func removeStopWords(tokens []string) (res []string) {
	res = []string{}
	for _, token := range tokens {
		if _, found := stopWords[token]; !found {
			res = append(res, token)
		}
	}
	return
}
