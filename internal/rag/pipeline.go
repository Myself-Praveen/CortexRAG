package rag

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/Myself-Praveen/CortexRAG/internal/document"
	"github.com/Myself-Praveen/CortexRAG/internal/embedding"
	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
	"github.com/Myself-Praveen/CortexRAG/internal/workerpool"
)

// Pipeline manages the ingestion of documents.
type Pipeline struct {
	mdExtractor  document.Extractor
	pdfExtractor document.Extractor
	chunker      *document.Chunker
	embedder     embedding.Provider
	store        vectorstore.Store
	pool         *workerpool.Pool
	cfg          *config.Config
}

func NewPipeline(
	embedder embedding.Provider,
	store vectorstore.Store,
	pool *workerpool.Pool,
	cfg *config.Config,
) *Pipeline {
	return &Pipeline{
		mdExtractor:  document.NewMarkdownExtractor(),
		pdfExtractor: document.NewPDFExtractor(),
		chunker:      document.NewChunker(),
		embedder:     embedder,
		store:        store,
		pool:         pool,
		cfg:          cfg,
	}
}

// Ingest extracts, chunks, embeds, and stores documents from the given files.
func (p *Pipeline) Ingest(ctx context.Context, files []string) error {
	docCh := make(chan *document.Document, len(files))
	chunkCh := make(chan *document.Chunk, 1000)

	var extractWg sync.WaitGroup

	// Stage 1: Extract (fan-out across files)
	for _, f := range files {
		extractWg.Add(1)
		go func(path string) {
			defer extractWg.Done()
			
			var doc *document.Document
			var err error
			
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".md":
				doc, err = p.mdExtractor.Extract(ctx, path)
			case ".pdf":
				doc, err = p.pdfExtractor.Extract(ctx, path)
			default:
				return // unsupported format
			}
			
			if err != nil {
				return // skip on error
			}
			docCh <- doc
		}(f)
	}

	go func() {
		extractWg.Wait()
		close(docCh)
	}()

	var chunkWg sync.WaitGroup

	// Stage 2: Chunk
	chunkWg.Add(1)
	go func() {
		defer chunkWg.Done()
		for doc := range docCh {
			chunks := p.chunker.Chunk(doc, p.cfg.RAG.ChunkSize, p.cfg.RAG.ChunkOverlap)
			for _, c := range chunks {
				chunkCh <- c
			}
		}
		close(chunkCh)
	}()

	var embedWg sync.WaitGroup

	// Stage 3: Embed + Store (worker pool)
	// We submit batches of chunks or individual chunks to the pool.
	// For simplicity, we submit one task per chunk.
	
	// Start a goroutine to read from chunkCh and submit tasks
	embedWg.Add(1)
	go func() {
		defer embedWg.Done()
		var taskWg sync.WaitGroup
		for chunk := range chunkCh {
			c := chunk // capture loop variable
			taskWg.Add(1)
			p.pool.Submit(func() {
				defer taskWg.Done()
				// Embed
				vec, err := p.embedder.Embed(ctx, c.Content)
				if err != nil {
					return
				}
				c.Embedding = vec
				
				// Store
				p.store.Upsert(ctx, c)
			})
		}
		taskWg.Wait()
	}()

	// Wait for pipeline completion
	chunkWg.Wait()
	embedWg.Wait()

	return nil
}
