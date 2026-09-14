package vectorstore

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	a := []float32{1, 0, 0}
	b := []float32{0, 1, 0}
	
	sim := CosineSimilarity(a, b)
	assert.Equal(t, float32(0.0), sim)

	c := []float32{1, 0, 0}
	sim2 := CosineSimilarity(a, c)
	assert.Equal(t, float32(1.0), sim2)

	// Edge case: empty vectors
	assert.Equal(t, float32(0.0), CosineSimilarity([]float32{}, []float32{}))
	
	// Edge case: zero vectors
	assert.Equal(t, float32(0.0), CosineSimilarity([]float32{0,0}, []float32{0,0}))
}

func BenchmarkCosineSimilarity(b *testing.B) {
	dim := 768
	v1 := make([]float32, dim)
	v2 := make([]float32, dim)
	
	for i := 0; i < dim; i++ {
		v1[i] = rand.Float32()
		v2[i] = rand.Float32()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CosineSimilarity(v1, v2)
	}
}
