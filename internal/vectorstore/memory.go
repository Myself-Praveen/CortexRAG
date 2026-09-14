package vectorstore

import (
	"context"
	"runtime"
	"sort"
	"sync"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
)

// StoredVector combines a chunk ID with its vector and original chunk.
type StoredVector struct {
	ID        string
	Embedding []float32
	Chunk     *document.Chunk
}

// MemoryStore is an in-memory brute-force vector store.
type MemoryStore struct {
	mu      sync.RWMutex
	vectors map[string]StoredVector
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		vectors: make(map[string]StoredVector),
	}
}

func (s *MemoryStore) Upsert(ctx context.Context, chunk *document.Chunk) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors[chunk.ID] = StoredVector{
		ID:        chunk.ID,
		Embedding: chunk.Embedding,
		Chunk:     chunk,
	}
	return nil
}

func (s *MemoryStore) UpsertBatch(ctx context.Context, chunks []*document.Chunk) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, chunk := range chunks {
		s.vectors[chunk.ID] = StoredVector{
			ID:        chunk.ID,
			Embedding: chunk.Embedding,
			Chunk:     chunk,
		}
	}
	return nil
}

func (s *MemoryStore) Delete(ctx context.Context, chunkID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vectors, chunkID)
	return nil
}

func (s *MemoryStore) DeleteDocument(ctx context.Context, docID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Naive approach: iterate and delete
	for id, vec := range s.vectors {
		if vec.Chunk.DocID == docID {
			delete(s.vectors, id)
		}
	}
	return nil
}

func (s *MemoryStore) Stats(ctx context.Context) (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]any{
		"count": len(s.vectors),
		"type":  "memory",
	}, nil
}

// Search performs a concurrent brute-force cosine similarity search using a worker pool concept inline.
func (s *MemoryStore) Search(ctx context.Context, query []float32, topK int) ([]SearchResult, error) {
	s.mu.RLock()
	allVectors := make([]StoredVector, 0, len(s.vectors))
	for _, v := range s.vectors {
		allVectors = append(allVectors, v)
	}
	s.mu.RUnlock()

	if len(allVectors) == 0 {
		return nil, nil
	}

	numWorkers := runtime.NumCPU()
	if len(allVectors) < numWorkers {
		numWorkers = 1
	}
	chunkSize := (len(allVectors) + numWorkers - 1) / numWorkers

	resultsCh := make(chan []SearchResult, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(allVectors) {
			end = len(allVectors)
		}

		go func(slice []StoredVector) {
			localResults := make([]SearchResult, 0, topK+1)
			for _, v := range slice {
				score := CosineSimilarity(query, v.Embedding)
				localResults = append(localResults, SearchResult{
					ChunkID: v.ID,
					Score:   score,
					Chunk:   v.Chunk,
				})
				
				// Keep sorted
				sort.Slice(localResults, func(i, j int) bool {
					return localResults[i].Score > localResults[j].Score
				})
				if len(localResults) > topK {
					localResults = localResults[:topK]
				}
			}
			resultsCh <- localResults
		}(allVectors[start:end])
	}

	merged := make([]SearchResult, 0, topK*numWorkers)
	for i := 0; i < numWorkers; i++ {
		partial := <-resultsCh
		merged = append(merged, partial...)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Score > merged[j].Score
	})

	if len(merged) > topK {
		merged = merged[:topK]
	}

	return merged, nil
}
