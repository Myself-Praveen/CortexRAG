package llm

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeminiProvider(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("Skipping Gemini test because GEMINI_API_KEY is not set")
	}

	provider, err := NewGeminiProvider(context.Background(), "gemini-1.5-flash")
	assert.NoError(t, err)
	assert.NotNil(t, provider)

	// Since we mock API without key, we can't reliably test Generate here without quota usage.
	// But structure is validated.
}
