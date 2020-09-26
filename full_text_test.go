package fts

import (
	"fts/corpus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TokenHas(tokens []string, token string) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}

func TestTokens(t *testing.T) {
	searcher := FullTextSearch{}
	doc := &corpus.Document{
		Text: "hello, I am vic.",
	}
	tokens := searcher.Tokenize(doc.Text)
	assert.Len(t, tokens, 4)
	assert.True(t, TokenHas(tokens, "hello"))
	assert.True(t, TokenHas(tokens, "I"))
	assert.True(t, TokenHas(tokens, "am"))
	assert.True(t, TokenHas(tokens, "vic"))
}

func TestSearch(t *testing.T) {
	searcher := NewFullText(docs)
	subDocs, err := searcher.Search(docs, "hello cat")
	assert.NoError(t, err)
	assert.Greater(t, subDocs.Len(), 0)
}
