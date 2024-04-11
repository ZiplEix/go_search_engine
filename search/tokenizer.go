package search

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

func analyze(text string) []string {
	token := tokenize(text)
	tokens := lowercaseFilter(token)
	tokens = stopWordFilter(token)
	tokens = stemmerFilter(token)

	return tokens
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLower(r) && !unicode.IsNumber(r)
	})
}

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}

	return r
}

func stopWordFilter(tokens []string) []string {
	var stopwords = map[string]struct{}{
		"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
		"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
		"it": {}, "for": {}, "not": {}, "on": {}, "with": {},
		"as": {}, "you": {}, "do": {}, "at": {}, "this": {},
		"but": {}, "his": {}, "by": {}, "from": {}, "they": {},
		"we": {}, "say": {}, "her": {}, "she": {}, "or": {},
		"an": {}, "will": {}, "my": {}, "one": {}, "all": {},
		"www": {}, "com": {}, "org": {}, "net": {}, "io": {},
		"https": {}, "http": {}, "html": {}, "php": {}, "asp": {}, "co": {}, "fr": {},
	}

	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}

	return r
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}

	return r
}
