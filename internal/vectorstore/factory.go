package vectorstore

import (
	"fmt"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
)

// Factory creates a vector store based on the configuration.
func Factory(cfg *config.Config) (Store, error) {
	// For simplicity, we just check the config
	// Since VectorDB doesn't exist, we fallback
	switch "memory" {
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
