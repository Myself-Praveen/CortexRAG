package embedding

import (
	"context"

	"google.golang.org/genai"
)

type GeminiEmbedder struct {
	client *genai.Client
	model  string
}

func NewGeminiEmbedder(ctx context.Context, model string) (*GeminiEmbedder, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	if model == "" {
		model = "gemini-embedding-001"
	}

	return &GeminiEmbedder{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	result, err := g.client.Models.EmbedContent(ctx, g.model,
		genai.Text(text),
		&genai.EmbedContentConfig{TaskType: "RETRIEVAL_DOCUMENT"},
	)
	if err != nil {
		return nil, err
	}
	return result.Embeddings[0].Values, nil
}

func (g *GeminiEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	result, err := g.client.Models.EmbedContent(ctx, g.model,
		genai.Text(text),
		&genai.EmbedContentConfig{TaskType: "RETRIEVAL_QUERY"},
	)
	if err != nil {
		return nil, err
	}
	return result.Embeddings[0].Values, nil
}

func (g *GeminiEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var results [][]float32
	// For simplicity, simulating batch as repeated calls here.
	// Gemini might support batch endpoints, but we do naive iteration.
	for _, t := range texts {
		vec, err := g.Embed(ctx, t)
		if err != nil {
			return nil, err
		}
		results = append(results, vec)
	}
	return results, nil
}
