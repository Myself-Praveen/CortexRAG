package llm

import "context"

// Provider defines the interface for language models.
type Provider interface {
	// Generate text from a prompt.
	Generate(ctx context.Context, prompt string) (string, error)
	// GenerateStream streams text from a prompt using channels.
	GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error)
}
