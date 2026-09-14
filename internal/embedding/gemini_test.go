package embedding

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeminiEmbedder_Mock(t *testing.T) {
	// Without an API key, we can only test the creation failure or skip actual API call
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("Skipping Gemini test because GEMINI_API_KEY is not set")
	}

	embedder, err := NewGeminiEmbedder(context.Background(), "gemini-embedding-001")
	assert.NoError(t, err)
	assert.NotNil(t, embedder)

	vec, err := embedder.Embed(context.Background(), "test")
	assert.NoError(t, err)
	assert.NotEmpty(t, vec)
}
