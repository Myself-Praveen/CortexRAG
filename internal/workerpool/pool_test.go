package workerpool

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	pool := New(5)
	
	var counter atomic.Int32
	
	for i := 0; i < 100; i++ {
		pool.Submit(func() {
			counter.Add(1)
		})
	}
	
	pool.Shutdown()
	
	assert.Equal(t, int32(100), counter.Load())
}
