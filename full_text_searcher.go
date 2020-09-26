package fts

import (
	"fts/corpus"
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

var stopwords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func NewFullText(docs *corpus.Documents) Searcher {
	search := &FullTextSearch{
		docs:    docs,
		docsMap: docs.ToMap(),
		index:   make(map[string][]int),
	}
	search.buildIndex(docs)
	return search
}

type FullTextSearch struct {
	docs    *corpus.Documents
	index   map[string][]int
	docsMap map[int]*corpus.Document
}

func (s *FullTextSearch) Search(docs *corpus.Documents, term string) (*corpus.Documents, error) {
	tokens := s.Normailze(s.Tokenize(term))

	ids := make([]int, 0, 1)

	for _, tok := range tokens {
		subIDs := s.index[tok]
		if len(ids) == 0 {
			ids = subIDs
		} else {
			ids = s.intersection(ids, subIDs)
		}
	}

	return s.findDocs(ids), nil
}

func (s *FullTextSearch) findDocs(ids []int) *corpus.Documents {
	docs := corpus.New()

	for _, id := range ids {
		docs.Add(s.docsMap[id])
	}
	return docs
}

func (s *FullTextSearch) Tokenize(term string) []string {
	return strings.FieldsFunc(term, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func (s *FullTextSearch) Normailze(tokens []string) []string {
	return s.stemming(s.stopwords(s.lowercase(tokens)))
}

func (s *FullTextSearch) intersection(a, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	result := make([]int, 0, maxLen)
	exist := make(map[int]bool)

	for _, bID := range b {
		exist[bID] = true
	}

	for _, aID := range a {
		if _, ok := exist[aID]; ok {
			result = append(result, aID)
		}
	}

	return result
}

func (s *FullTextSearch) buildIndex(docs *corpus.Documents) error {
	index := make(map[string][]int)

	ranger := docs.Iterator()

	for {
		_, doc, next := ranger()
		if !next {
			break
		}
		tokens := s.Normailze(s.Tokenize(doc.Text))

		for _, t := range tokens {
			ids := index[t]
			if ids == nil || cap(ids) == 0 {
				ids = make([]int, 0, 1)
			}
			if len(ids) > 0 && ids[len(ids)-1] == doc.ID {
				continue
			}
			index[t] = append(ids, doc.ID)
		}
	}
	s.index = index
	return nil
}

func (s *FullTextSearch) lowercase(tokens []string) []string {
	result := make([]string, len(tokens))

	for i := range tokens {
		result[i] = strings.ToLower(tokens[i])
	}

	return result
}

func (s *FullTextSearch) stopwords(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			result = append(result, token)
		}
	}
	return result
}

func (s *FullTextSearch) stemming(tokens []string) []string {
	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = snowballeng.Stem(token, false)
	}
	return result
}
