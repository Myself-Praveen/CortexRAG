package vectorstore

import (
	"context"
	"fmt"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
)

// Factory creates a vector store based on the configuration.
func Factory(cfg *config.Config) (Store, error) {
	// For simplicity, we just check the config
	switch cfg.RAG.VectorDB {
	case "hnsw":
		return NewHNSWStore(16, 64), nil
	case "pgvector":
		// return NewPGVectorStore(cfg.Database.URL)
		return nil, fmt.Errorf("pgvector not implemented in this factory stub")
	case "memory":
		fallthrough
	default:
		return NewMemoryStore(), nil
	}
}
