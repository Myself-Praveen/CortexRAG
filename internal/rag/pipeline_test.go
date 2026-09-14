package rag

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/Myself-Praveen/CortexRAG/internal/document"
	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
	"github.com/Myself-Praveen/CortexRAG/internal/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockProvider
type mockEmbedder struct{}
func (m *mockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	return []float32{1.0, 0.0}, nil
}
func (m *mockEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return []float32{1.0, 0.0}, nil
}
func (m *mockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, nil
}

// MockStore
type mockStore struct {
	chunks []*document.Chunk
}
func (m *mockStore) Upsert(ctx context.Context, chunk *document.Chunk) error {
	m.chunks = append(m.chunks, chunk)
	return nil
}
func (m *mockStore) UpsertBatch(ctx context.Context, chunks []*document.Chunk) error { return nil }
func (m *mockStore) Search(ctx context.Context, query []float32, topK int) ([]vectorstore.SearchResult, error) { return nil, nil }
func (m *mockStore) Delete(ctx context.Context, chunkID string) error { return nil }
func (m *mockStore) DeleteDocument(ctx context.Context, docID string) error { return nil }
func (m *mockStore) Stats(ctx context.Context) (map[string]any, error) { return nil, nil }


func TestPipeline(t *testing.T) {
	pool := workerpool.New(2)
	defer pool.Shutdown()

	embedder := &mockEmbedder{}
	store := &mockStore{}
	
	cfg := &config.Config{}
	cfg.RAG.ChunkSize = 10
	cfg.RAG.ChunkOverlap = 2

	pipeline := NewPipeline(embedder, store, pool, cfg)

	// Setup dummy file
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10"
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pipeline.Ingest(ctx, []string{path})
	require.NoError(t, err)
	
	assert.Greater(t, len(store.chunks), 0)
}
