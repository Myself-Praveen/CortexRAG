package vectorstore

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
)

func BenchmarkStores(b *testing.B) {
	ctx := context.Background()
	dim := 768
	numChunks := 1000

	memStore := NewMemoryStore()
	hnswStore := NewHNSWStore(16, 64)

	// Pre-populate
	for i := 0; i < numChunks; i++ {
		vec := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[j] = rand.Float32()
		}
		
		chunk := &document.Chunk{
			ID:        fmt.Sprintf("chunk-%d", i),
			Embedding: vec,
		}

		_ = memStore.Upsert(ctx, chunk)
		_ = hnswStore.Upsert(ctx, chunk)
	}

	query := make([]float32, dim)
	for j := 0; j < dim; j++ {
		query[j] = rand.Float32()
	}

	b.Run("MemoryStore_BruteForce", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = memStore.Search(ctx, query, 10)
		}
	})

	b.Run("HNSWStore_GraphSearch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = hnswStore.Search(ctx, query, 10)
		}
	})
}
