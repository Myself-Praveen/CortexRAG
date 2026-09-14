package vectorstore

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPGVectorStore(t *testing.T) {
	t.Skip("Skipping pgvector integration tests (requires docker/postgres)")
	
	store, err := NewPGVectorStore("postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	assert.NoError(t, err)
	
	stats, err := store.Stats(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "pgvector", stats["type"])
}
