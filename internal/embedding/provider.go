package embedding

import "context"

// Provider defines the interface for interacting with embedding models.
type Provider interface {
	// Embed generates an embedding for a document chunk.
	Embed(ctx context.Context, text string) ([]float32, error)
	// EmbedQuery generates an embedding for a search query.
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
	// EmbedBatch generates embeddings for multiple chunks concurrently.
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}
