package rag

import (
	"context"
	"testing"
	"time"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/stretchr/testify/assert"
)

type mockLLM struct{}
func (m *mockLLM) Generate(ctx context.Context, prompt string) (string, error) {
	return "Mocked Answer", nil
}
func (m *mockLLM) GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	c := make(chan string, 1)
	e := make(chan error, 1)
	c <- "Mocked "
	c <- "Answer"
	close(c)
	close(e)
	return c, e
}

func TestGenerator(t *testing.T) {
	embedder := &mockEmbedder{} // from pipeline_test.go
	store := &mockStore{}
	llm := &mockLLM{}
	
	cfg := &config.Config{}
	cfg.RAG.TopK = 3

	gen := NewGenerator(embedder, store, llm, cfg)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ans, err := gen.Answer(ctx, "test query")
	assert.NoError(t, err)
	assert.Equal(t, "Mocked Answer", ans)

	chunks, errs := gen.AnswerStream(ctx, "test query")
	var result string
	for c := range chunks {
		result += c
	}
	
	for err := range errs {
		assert.NoError(t, err)
	}

	assert.Equal(t, "Mocked Answer", result)
}
