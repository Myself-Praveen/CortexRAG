package embedding

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOllamaEmbedder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"embedding": [0.1, 0.2, 0.3]}`))
	}))
	defer ts.Close()

	embedder := NewOllamaEmbedder(ts.URL, "test-model")

	vec, err := embedder.Embed(context.Background(), "test text")
	assert.NoError(t, err)
	assert.Equal(t, []float32{0.1, 0.2, 0.3}, vec)
}
