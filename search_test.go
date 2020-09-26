package fts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSimple(b *testing.B) {
	searcher := New(Simple)
	catDocs, err := searcher.Search(docs, "cat")
	assert.NoError(b, err)
	assert.Greater(b, catDocs.Len(), 1)
}

func BenchmarkRegex(b *testing.B) {
	searcher := New(Regex)
	catDocs, err := searcher.Search(docs, "cat")
	assert.NoError(b, err)
	assert.Greater(b, catDocs.Len(), 1)
}

func BenchmarkInvertIndex(b *testing.B) {
	catDocs, err := fullTextSearch.Search(docs, "cat")
	assert.NoError(b, err)
	assert.Greater(b, catDocs.Len(), 1)
}
