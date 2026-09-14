package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOllamaProvider_Generate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response": "test response", "done": true}`))
	}))
	defer ts.Close()

	provider := NewOllamaProvider(ts.URL, "test-model")

	resp, err := provider.Generate(context.Background(), "test prompt")
	assert.NoError(t, err)
	assert.Equal(t, "test response", resp)
}

func TestOllamaProvider_GenerateStream(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response": "part 1", "done": false}` + "\n" + `{"response": " part 2", "done": true}`))
	}))
	defer ts.Close()

	provider := NewOllamaProvider(ts.URL, "test-model")

	chunks, errs := provider.GenerateStream(context.Background(), "test prompt")
	
	var result string
	for c := range chunks {
		result += c
	}
	
	for err := range errs {
		assert.NoError(t, err)
	}

	assert.Equal(t, "part 1 part 2", result)
}
