package searchengine

import (
	"regexp"
	"strings"
)

// punctuationRegex is used to remove all non-word and non-space characters from text.
var punctuationRegex = regexp.MustCompile(`[^\w\s]`)

// normalize converts text to lowercase and removes punctuation.
func normalize(text string) string {
	text = strings.ToLower(text)
	text = punctuationRegex.ReplaceAllString(text, "")
	return text
}

// tokenize splits the normalized text into individual tokens/words.
func tokenize(text string) (res []string) {
	res = strings.Fields(text)
	return
}