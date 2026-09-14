package vectorstore

import (
	"context"
	"fmt"
	"testing"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
	"github.com/stretchr/testify/assert"
)

func TestHNSWStore(t *testing.T) {
	store := NewHNSWStore(16, 64)
	ctx := context.Background()

	// Insert
	for i := 0; i < 20; i++ {
		err := store.Upsert(ctx, &document.Chunk{
			ID:        fmt.Sprintf("chunk-%d", i),
			DocID:     "doc-1",
			Embedding: []float32{float32(i), 0, 0},
			Content:   fmt.Sprintf("content %d", i),
		})
		assert.NoError(t, err)
	}

	stats, _ := store.Stats(ctx)
	assert.Equal(t, 20, stats["count"])

	// Search
	query := []float32{1, 0, 0}
	results, err := store.Search(ctx, query, 5)
	assert.NoError(t, err)
	assert.Len(t, results, 5) 
}
