package embedding

import (
	"context"
	"sync"
)

// BatchEmbedder runs embeddings concurrently if the underlying provider doesn't support batching natively.
type BatchEmbedder struct {
	provider Provider
}

func NewBatchEmbedder(provider Provider) *BatchEmbedder {
	return &BatchEmbedder{provider: provider}
}

func (b *BatchEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	return b.provider.Embed(ctx, text)
}

func (b *BatchEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return b.provider.EmbedQuery(ctx, text)
}

func (b *BatchEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	results := make([][]float32, len(texts))
	errs := make([]error, len(texts))

	var wg sync.WaitGroup
	
	// Semaphore to limit concurrency
	sem := make(chan struct{}, 10)

	for i, text := range texts {
		wg.Add(1)
		go func(index int, t string) {
			defer wg.Done()
			
			sem <- struct{}{}
			defer func() { <-sem }()
			
			vec, err := b.provider.Embed(ctx, t)
			results[index] = vec
			errs[index] = err
		}(i, text)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
