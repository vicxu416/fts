package fts

import (
	"fts/corpus"
	"log"
	"regexp"
	"strings"
)

var (
	docs           *corpus.Documents
	fullTextSearch Searcher
)

func init() {
	var err error
	docs, err = corpus.GetAll()
	if err != nil {
		log.Panicf("get documents failed, err:%+v", err)
	}
	log.Println("get documents successfully")
	fullTextSearch = NewFullText(docs)
}

func New(way SearchWay) Searcher {
	switch way {
	case Simple:
		return SearchFunc(SimplSearch)
	case Regex:
		return SearchFunc(RegexSearch)
	}
	return nil
}

type SearchWay int8

const (
	Simple SearchWay = iota + 1
	Regex
	FullText
)

type Searcher interface {
	Search(docs *corpus.Documents, term string) (*corpus.Documents, error)
}

type SearchFunc func(docs *corpus.Documents, term string) (*corpus.Documents, error)

func (s SearchFunc) Search(docs *corpus.Documents, term string) (*corpus.Documents, error) {
	return s(docs, term)
}

func SimplSearch(docs *corpus.Documents, term string) (*corpus.Documents, error) {
	result := corpus.New()
	ranger := docs.Iterator()

	for {
		_, doc, next := ranger()
		if !next {
			break
		}
		if strings.Contains(doc.Text, term) {
			result.Add(doc)
		}
	}

	return result, nil
}

func RegexSearch(docs *corpus.Documents, term string) (*corpus.Documents, error) {
	re, err := regexp.Compile(`(?i)\b` + term + `\b`)
	if err != nil {
		return nil, err
	}
	result := corpus.New()

	ranger := docs.Iterator()

	for {
		_, doc, next := ranger()
		if !next {
			break
		}

		if re.MatchString(doc.Text) {
			result.Add(doc)
		}
	}
	return result, nil
}
