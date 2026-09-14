package rag

import (
	"context"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/Myself-Praveen/CortexRAG/internal/embedding"
	"github.com/Myself-Praveen/CortexRAG/internal/llm"
	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
)

type Generator struct {
	embedder embedding.Provider
	store    vectorstore.Store
	llm      llm.Provider
	cfg      *config.Config
}

func NewGenerator(
	embedder embedding.Provider,
	store vectorstore.Store,
	llmProvider llm.Provider,
	cfg *config.Config,
) *Generator {
	return &Generator{
		embedder: embedder,
		store:    store,
		llm:      llmProvider,
		cfg:      cfg,
	}
}

func (g *Generator) Answer(ctx context.Context, question string) (string, error) {
	// 1. Embed Question
	queryVec, err := g.embedder.EmbedQuery(ctx, question)
	if err != nil {
		return "", err
	}

	// 2. Search Store
	results, err := g.store.Search(ctx, queryVec, g.cfg.RAG.TopK)
	if err != nil {
		return "", err
	}

	// 3. Build Prompt
	prompt := llm.BuildPrompt(question, results)

	// 4. Generate Answer
	return g.llm.Generate(ctx, prompt)
}

func (g *Generator) AnswerStream(ctx context.Context, question string) (<-chan string, <-chan error) {
	// For streaming, we do steps 1-3 synchronously (or in a goroutine), then stream.
	// We'll wrap the setup in a goroutine so it doesn't block the caller.
	
	chunks := make(chan string)
	errs := make(chan error, 1)

	go func() {
		// 1. Embed
		queryVec, err := g.embedder.EmbedQuery(ctx, question)
		if err != nil {
			errs <- err
			close(chunks)
			close(errs)
			return
		}

		// 2. Search
		results, err := g.store.Search(ctx, queryVec, g.cfg.RAG.TopK)
		if err != nil {
			errs <- err
			close(chunks)
			close(errs)
			return
		}

		// 3. Build Prompt
		prompt := llm.BuildPrompt(question, results)

		// 4. Generate Stream
		llmChunks, llmErrs := g.llm.GenerateStream(ctx, prompt)
		
		for {
			select {
			case c, ok := <-llmChunks:
				if !ok {
					llmChunks = nil
				} else {
					chunks <- c
				}
			case e, ok := <-llmErrs:
				if !ok {
					llmErrs = nil
				} else if e != nil {
					errs <- e
				}
			}

			if llmChunks == nil && llmErrs == nil {
				break
			}
		}

		close(chunks)
		close(errs)
	}()

	return chunks, errs
}
