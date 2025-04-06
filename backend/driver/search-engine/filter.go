package searchengine

var stopWords = map[string]struct{}{
	"the": {}, "is": {}, "and": {}, "a": {}, "an": {}, "in": {}, "of": {}, "on": {}, "at": {}, "to": {},
}

func removeStopWords(tokens []string) (res []string) {
	res = []string{}
	for _, token := range tokens {
		if _, found := stopWords[token]; !found {
			res = append(res, token)
		}
	}
	return
}
