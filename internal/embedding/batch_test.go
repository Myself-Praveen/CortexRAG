package embedding

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchEmbedder(t *testing.T) {
	provider := &dummyProvider{}
	batcher := NewBatchEmbedder(provider)

	texts := []string{"one", "two", "three", "four", "five"}
	results, err := batcher.EmbedBatch(context.Background(), texts)
	
	assert.NoError(t, err)
	assert.Len(t, results, 5)
	
	for _, res := range results {
		assert.Equal(t, []float32{1, 2, 3}, res)
	}
}
