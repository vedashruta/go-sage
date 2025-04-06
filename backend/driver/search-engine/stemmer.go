package searchengine

import "github.com/kljensen/snowball"

// stemTokens applies stemming to each token using the Snowball stemmer.
// Only successful stems are included in the result.
func stemTokens(tokens []string) (res []string) {
	res = []string{}
	for _, token := range tokens {
		stem, err := snowball.Stem(token, "english", true)
		if err == nil {
			res = append(res, stem)
		}
	}
	return
}
