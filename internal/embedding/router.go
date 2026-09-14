package embedding

import (
	"context"
	"fmt"
)

type Router struct {
	providers []Provider
}

// NewRouter creates a router that tries providers in order.
func NewRouter(providers ...Provider) *Router {
	return &Router{
		providers: providers,
	}
}

func (r *Router) Embed(ctx context.Context, text string) ([]float32, error) {
	var lastErr error
	for _, p := range r.providers {
		vec, err := p.Embed(ctx, text)
		if err == nil {
			return vec, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all embedding providers failed, last error: %w", lastErr)
}

func (r *Router) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	var lastErr error
	for _, p := range r.providers {
		vec, err := p.EmbedQuery(ctx, text)
		if err == nil {
			return vec, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all embedding providers failed on query, last error: %w", lastErr)
}

func (r *Router) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var lastErr error
	for _, p := range r.providers {
		vec, err := p.EmbedBatch(ctx, texts)
		if err == nil {
			return vec, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all embedding providers failed on batch, last error: %w", lastErr)
}
