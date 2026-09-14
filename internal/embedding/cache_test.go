package embedding

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyProvider struct {
	calls int
	mu    sync.Mutex
}

func (d *dummyProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	d.mu.Lock()
	d.calls++
	d.mu.Unlock()
	return []float32{1, 2, 3}, nil
}
func (d *dummyProvider) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return d.Embed(ctx, text)
}
func (d *dummyProvider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, nil
}

func TestCachedEmbedder(t *testing.T) {
	provider := &dummyProvider{}
	cache := NewCachedEmbedder(provider, 32) // small capacity

	t.Run("Cache Miss then Hit", func(t *testing.T) {
		ctx := context.Background()
		
		vec1, err := cache.Embed(ctx, "test query")
		assert.NoError(t, err)
		assert.Equal(t, []float32{1, 2, 3}, vec1)
		
		vec2, err := cache.Embed(ctx, "test query")
		assert.NoError(t, err)
		assert.Equal(t, []float32{1, 2, 3}, vec2)
		
		assert.Equal(t, 1, provider.calls)
		
		hits, misses := cache.Stats()
		assert.Equal(t, int64(1), hits)
		assert.Equal(t, int64(1), misses)
	})

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				_, _ = cache.Embed(context.Background(), fmt.Sprintf("query %d", val%10)) // 10 unique queries
			}(i)
		}
		wg.Wait()
		
		// Wait... provider.calls was 1 before this. Now we sent 10 unique queries.
		// So total calls should be 11.
		assert.Equal(t, 11, provider.calls)
	})
}
