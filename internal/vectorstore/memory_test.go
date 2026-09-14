package vectorstore

import (
	"context"
	"fmt"
	"testing"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Insert
	for i := 0; i < 10; i++ {
		err := store.Upsert(ctx, &document.Chunk{
			ID:        fmt.Sprintf("chunk-%d", i),
			DocID:     "doc-1",
			Embedding: []float32{float32(i), 0, 0},
			Content:   fmt.Sprintf("content %d", i),
		})
		assert.NoError(t, err)
	}

	stats, _ := store.Stats(ctx)
	assert.Equal(t, 10, stats["count"])

	// Search
	query := []float32{1, 0, 0}
	results, err := store.Search(ctx, query, 3)
	assert.NoError(t, err)
	assert.Len(t, results, 3) // It should return top 3
	
	// They all have {something, 0, 0}, and query is {1, 0, 0}.
	// Except chunk-0 which is {0,0,0}. Cosine similarity for {x,0,0} with {1,0,0} is 1.0 (if x>0).
	// So many will have score 1.0

	// Delete
	err = store.Delete(ctx, "chunk-1")
	assert.NoError(t, err)
	
	stats, _ = store.Stats(ctx)
	assert.Equal(t, 9, stats["count"])

	// Delete Document
	err = store.DeleteDocument(ctx, "doc-1")
	assert.NoError(t, err)

	stats, _ = store.Stats(ctx)
	assert.Equal(t, 0, stats["count"])
}
