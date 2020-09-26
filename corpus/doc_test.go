package corpus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	docs, err := GetAll()
	assert.NoError(t, err)
	assert.Greater(t, docs.Len(), 0)
	t.Logf("Document Size:%d", docs.Len())
	assert.NotZero(t, docs.Docs[1].ID)
}
