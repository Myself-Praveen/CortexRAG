package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaEmbedder struct {
	url    string
	model  string
	client *http.Client
}

func NewOllamaEmbedder(url, model string) *OllamaEmbedder {
	if url == "" {
		url = "http://localhost:11434"
	}
	if model == "" {
		model = "nomic-embed-text"
	}
	return &OllamaEmbedder{
		url:    url,
		model:  model,
		client: &http.Client{},
	}
}

type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (o *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	reqBody, _ := json.Marshal(ollamaEmbedRequest{
		Model:  o.model,
		Prompt: text,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", o.url+"/api/embeddings", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API error: %d - %s", resp.StatusCode, string(body))
	}

	var res ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Embedding, nil
}

func (o *OllamaEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return o.Embed(ctx, text)
}

func (o *OllamaEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var results [][]float32
	for _, t := range texts {
		vec, err := o.Embed(ctx, t)
		if err != nil {
			return nil, err
		}
		results = append(results, vec)
	}
	return results, nil
}
