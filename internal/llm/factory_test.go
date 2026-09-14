package llm

import (
	"context"
	"testing"
	"os"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	cfg := &config.Config{}
	cfg.LLM.Provider = "ollama"
	cfg.LLM.Model = "test-model"

	provider, err := Factory(context.Background(), cfg)
	assert.NoError(t, err)
	assert.IsType(t, &OllamaProvider{}, provider)

	cfg.LLM.Provider = "gemini"
	if os.Getenv("GEMINI_API_KEY") != "" {
		provider2, err := Factory(context.Background(), cfg)
		assert.NoError(t, err)
		assert.IsType(t, &GeminiProvider{}, provider2)
	}
}
