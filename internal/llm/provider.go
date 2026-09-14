package llm

import "context"

// Provider defines the interface for language models.
type Provider interface {
	// Generate text from a prompt.
	Generate(ctx context.Context, prompt string, systemPrompt string) (string, error)
	// StreamGenerate streams text from a prompt using a callback function.
	StreamGenerate(ctx context.Context, prompt string, systemPrompt string, onToken func(string)) error
}
