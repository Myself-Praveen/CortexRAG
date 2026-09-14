package llm

import (
	"context"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func NewGeminiProvider(ctx context.Context, model string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &GeminiProvider{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiProvider) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return resp.Text(), nil
	}
	return "", nil
}

func (g *GeminiProvider) GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	chunks := make(chan string, 50)
	errs := make(chan error, 1)

	go func() {
		defer close(chunks)
		defer close(errs)

		for resp, err := range g.client.Models.GenerateContentStream(ctx, g.model, genai.Text(prompt), nil) {
			if err != nil {
				errs <- err
				return
			}
			chunks <- resp.Text()
		}
	}()

	return chunks, errs
}
