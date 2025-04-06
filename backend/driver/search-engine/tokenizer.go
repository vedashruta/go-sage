package searchengine

import (
	"regexp"
	"strings"
)

var punctuationRegex = regexp.MustCompile(`[^\w\s]`)

func normalize(text string) string {
	text = strings.ToLower(text)
	text = punctuationRegex.ReplaceAllString(text, "")
	return text
}

func tokenize(text string) (res []string) {
	res = strings.Fields(text)
	return
}
