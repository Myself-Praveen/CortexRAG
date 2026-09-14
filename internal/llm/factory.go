package llm

import (
	"context"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
)

// Factory creates an LLM provider based on config.
func Factory(ctx context.Context, cfg *config.Config) (Provider, error) {
	switch cfg.LLM.Provider {
	case "ollama":
		return NewOllamaProvider(cfg.LLM.Model, cfg.LLM.Model), nil
	case "gemini":
		fallthrough
	default:
		return NewGeminiProvider(ctx, cfg.LLM.Model)
	}
}
