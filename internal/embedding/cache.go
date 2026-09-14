package embedding

import (
	"container/list"
	"context"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

const numShards = 16

type CachedEmbedder struct {
	provider Provider
	shards   [numShards]*cacheShard
	hits     atomic.Int64
	misses   atomic.Int64
}

type cacheShard struct {
	mu    sync.RWMutex
	items map[uint64]*list.Element
	order *list.List
	cap   int
}

type cacheEntry struct {
	key       uint64
	text      string
	embedding []float32
}

func NewCachedEmbedder(provider Provider, maxSize int) *CachedEmbedder {
	shardCap := maxSize / numShards
	if shardCap <= 0 {
		shardCap = 100
	}

	c := &CachedEmbedder{
		provider: provider,
	}

	for i := 0; i < numShards; i++ {
		c.shards[i] = &cacheShard{
			items: make(map[uint64]*list.Element),
			order: list.New(),
			cap:   shardCap,
		}
	}

	return c
}

func hashText(text string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(text))
	return h.Sum64()
}

func (c *CachedEmbedder) getShard(key uint64) *cacheShard {
	return c.shards[key%numShards]
}

func (c *CachedEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	key := hashText(text)
	shard := c.getShard(key)

	shard.mu.RLock()
	if elem, ok := shard.items[key]; ok {
		entry := elem.Value.(*cacheEntry)
		if entry.text == text { // Handle hash collisions
			shard.mu.RUnlock()
			
			// Move to front (requires write lock)
			shard.mu.Lock()
			// Double check it's still there
			if e, ok := shard.items[key]; ok {
				shard.order.MoveToFront(e)
			}
			shard.mu.Unlock()
			
			c.hits.Add(1)
			return entry.embedding, nil
		}
	}
	shard.mu.RUnlock()

	c.misses.Add(1)

	// Cache miss, call underlying provider
	vec, err := c.provider.Embed(ctx, text)
	if err != nil {
		return nil, err
	}

	shard.mu.Lock()
	if len(shard.items) >= shard.cap {
		// Evict oldest
		oldest := shard.order.Back()
		if oldest != nil {
			oldEntry := oldest.Value.(*cacheEntry)
			delete(shard.items, oldEntry.key)
			shard.order.Remove(oldest)
		}
	}
	
	entry := &cacheEntry{key: key, text: text, embedding: vec}
	elem := shard.order.PushFront(entry)
	shard.items[key] = elem
	shard.mu.Unlock()

	return vec, nil
}

func (c *CachedEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	// Usually queries are unique, caching might be less effective but we do it anyway
	return c.Embed(ctx, text)
}

func (c *CachedEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	// For a real implementation, we should check cache for all, then send the missing ones to provider.EmbedBatch
	var results [][]float32
	for _, t := range texts {
		vec, err := c.Embed(ctx, t)
		if err != nil {
			return nil, err
		}
		results = append(results, vec)
	}
	return results, nil
}

func (c *CachedEmbedder) Stats() (hits int64, misses int64) {
	return c.hits.Load(), c.misses.Load()
}
