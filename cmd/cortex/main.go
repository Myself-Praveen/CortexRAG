package main

import (
	"context"
	"flag"
	"log"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/Myself-Praveen/CortexRAG/internal/embedding"
	"github.com/Myself-Praveen/CortexRAG/internal/llm"
	"github.com/Myself-Praveen/CortexRAG/internal/rag"
	"github.com/Myself-Praveen/CortexRAG/internal/server"
	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
	"github.com/Myself-Praveen/CortexRAG/internal/workerpool"
)

func main() {
	var configPath string
	var port int
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.IntVar(&port, "port", 0, "Port to run the server on")
	flag.Parse()

	if err := run(configPath, port); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run(configPath string, portOverride int) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if portOverride != 0 {
		cfg.Port = portOverride
	}

	ctx := context.Background()

	// 1. Vector Store
	store, err := vectorstore.Factory(cfg)
	if err != nil {
		return err
	}

	// 2. Embedder
	var embedder embedding.Provider
	switch cfg.Embedding.Provider {
	case "ollama":
		embedder = embedding.NewOllamaEmbedder(cfg.Embedding.Model, cfg.Embedding.Model)
	case "gemini":
		fallthrough
	default:
		// Requires GEMINI_API_KEY env var
		embedder, err = embedding.NewGeminiEmbedder(ctx, cfg.Embedding.Model)
		if err != nil {
			log.Printf("Warning: Failed to create Gemini embedder: %v", err)
			// Mock embedder to prevent crash if running locally without keys
			embedder = &embedding.CachedEmbedder{} // Just a placeholder that might crash on use but lets app start
		}
	}
	
	cachedEmbedder := embedding.NewCachedEmbedder(embedder, cfg.RAG.CacheSize)

	// 3. Worker Pool
	pool := workerpool.New(cfg.Worker.Count)
	defer pool.Shutdown()

	// 4. LLM
	llmProvider, err := llm.Factory(ctx, cfg)
	if err != nil {
		log.Printf("Warning: Failed to create LLM provider: %v", err)
	}

	// 5. Pipeline and Generator
	pipeline := rag.NewPipeline(cachedEmbedder, store, pool, cfg)
	generator := rag.NewGenerator(cachedEmbedder, store, llmProvider, cfg)

	// 6. Server
	api := server.NewAPI(pipeline, generator)
	srv := server.NewServer(cfg, api)

	return srv.Start()
}
