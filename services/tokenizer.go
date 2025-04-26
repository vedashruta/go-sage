package services

import (
	"regexp"
	"strings"
)

var punctuationRegex = regexp.MustCompile(`[^\w\s]`)

func Normalize(text string) string {
	text = strings.ToLower(text)
	text = punctuationRegex.ReplaceAllString(text, "")
	return text
}

func Tokenize(text string) (res []string) {
	res = strings.Fields(text)
	return
}
