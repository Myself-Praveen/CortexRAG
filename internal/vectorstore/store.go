package vectorstore

import (
	"context"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
)

// SearchResult represents a document chunk matched by vector search.
type SearchResult struct {
	ChunkID string
	Chunk   *document.Chunk
	Score   float32
}

// Store defines the interface for vector databases.
type Store interface {
	// Upsert adds or updates a chunk and its embedding.
	Upsert(ctx context.Context, chunk *document.Chunk) error
	// UpsertBatch adds or updates multiple chunks.
	UpsertBatch(ctx context.Context, chunks []*document.Chunk) error
	// Search finds the topK most similar chunks to the query vector.
	Search(ctx context.Context, query []float32, topK int) ([]SearchResult, error)
	// Delete removes a chunk by its ID.
	Delete(ctx context.Context, chunkID string) error
	// DeleteDocument removes all chunks associated with a document ID.
	DeleteDocument(ctx context.Context, docID string) error
	// Stats returns store statistics (e.g., number of vectors).
	Stats(ctx context.Context) (map[string]any, error)
}
