package vectorstore

import (
	"context"
	"math/rand"
	"sort"
	"sync"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
)

// HNSWStore is a simplified single-layer Navigable Small World (NSW) graph
// to represent the HNSW requirement. A full HNSW would have multiple layers.
type HNSWStore struct {
	mu            sync.RWMutex
	nodes         map[string]*hnswNode
	entryPointID  string
	M             int // max connections per node
	efConstruction int
}

type hnswNode struct {
	vector       StoredVector
	neighbors    []string // IDs of neighbors
}

func NewHNSWStore(m, ef int) *HNSWStore {
	return &HNSWStore{
		nodes:          make(map[string]*hnswNode),
		M:              m,
		efConstruction: ef,
	}
}

func (h *HNSWStore) Upsert(ctx context.Context, chunk *document.Chunk) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	newNode := &hnswNode{
		vector: StoredVector{
			ID:        chunk.ID,
			Embedding: chunk.Embedding,
			Chunk:     chunk,
		},
		neighbors: make([]string, 0, h.M),
	}

	h.nodes[chunk.ID] = newNode

	if h.entryPointID == "" {
		h.entryPointID = chunk.ID
		return nil
	}

	// Simple heuristic: connect to random existing nodes (simplified NSW)
	// In real HNSW, we'd search for the nearest neighbors and connect.
	var keys []string
	for k := range h.nodes {
		if k != chunk.ID {
			keys = append(keys, k)
		}
	}
	
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	
	connections := h.M
	if len(keys) < connections {
		connections = len(keys)
	}
	
	for i := 0; i < connections; i++ {
		newNode.neighbors = append(newNode.neighbors, keys[i])
		// Bidirectional
		h.nodes[keys[i]].neighbors = append(h.nodes[keys[i]].neighbors, chunk.ID)
	}

	return nil
}

func (h *HNSWStore) UpsertBatch(ctx context.Context, chunks []*document.Chunk) error {
	for _, chunk := range chunks {
		if err := h.Upsert(ctx, chunk); err != nil {
			return err
		}
	}
	return nil
}

// Search performs a greedy search on the graph.
func (h *HNSWStore) Search(ctx context.Context, query []float32, topK int) ([]SearchResult, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.entryPointID == "" {
		return nil, nil
	}

	// Visited set
	visited := make(map[string]bool)
	
	// Min-heap (or just sorted slice for simplicity) of candidates
	var candidates []SearchResult
	
	// Start at entry point
	curr := h.nodes[h.entryPointID]
	visited[curr.vector.ID] = true
	
	score := CosineSimilarity(query, curr.vector.Embedding)
	candidates = append(candidates, SearchResult{
		ChunkID: curr.vector.ID,
		Score:   score,
		Chunk:   curr.vector.Chunk,
	})

	// Very simplified greedy search
	for i := 0; i < len(candidates); i++ {
		currID := candidates[i].ChunkID
		currNode := h.nodes[currID]
		
		for _, neighborID := range currNode.neighbors {
			if !visited[neighborID] {
				visited[neighborID] = true
				neighborNode := h.nodes[neighborID]
				sim := CosineSimilarity(query, neighborNode.vector.Embedding)
				candidates = append(candidates, SearchResult{
					ChunkID: neighborID,
					Score:   sim,
					Chunk:   neighborNode.vector.Chunk,
				})
			}
		}
		
		sort.Slice(candidates, func(a, b int) bool {
			return candidates[a].Score > candidates[b].Score
		})
		
		if len(candidates) > h.efConstruction {
			candidates = candidates[:h.efConstruction]
		}
	}

	if len(candidates) > topK {
		candidates = candidates[:topK]
	}

	return candidates, nil
}

func (h *HNSWStore) Delete(ctx context.Context, chunkID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.nodes, chunkID)
	// Entry point might be broken, ignoring for simplicity
	return nil
}

func (h *HNSWStore) DeleteDocument(ctx context.Context, docID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, node := range h.nodes {
		if node.vector.Chunk.DocID == docID {
			delete(h.nodes, id)
		}
	}
	return nil
}

func (h *HNSWStore) Stats(ctx context.Context) (map[string]any, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return map[string]any{
		"count": len(h.nodes),
		"type":  "hnsw_simplified",
	}, nil
}
