package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type failingProvider struct{}
func (f *failingProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	return nil, errors.New("failed")
}
func (f *failingProvider) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return nil, errors.New("failed")
}
func (f *failingProvider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, errors.New("failed")
}

func TestRouterFallback(t *testing.T) {
	p1 := &failingProvider{}
	p2 := &dummyProvider{} // from cache_test.go, returns [1, 2, 3]

	router := NewRouter(p1, p2)
	
	vec, err := router.Embed(context.Background(), "test")
	assert.NoError(t, err)
	assert.Equal(t, []float32{1, 2, 3}, vec)
}
