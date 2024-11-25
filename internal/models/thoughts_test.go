package models

import (
	"testing"

	"github.com/tiwanakd/mythoughts-go/internal/assert"
)

func TestAddLikeandDislike(t *testing.T) {
	db := newTestDB(t)

	m := ThoughtModel{DB: db}

	likes, err := m.AddLike(1)
	assert.Equal(t, likes, 1)
	assert.NilError(t, err)

	dislikes, err := m.AddDislike(1)
	assert.Equal(t, dislikes, 1)
	assert.NilError(t, err)
}
