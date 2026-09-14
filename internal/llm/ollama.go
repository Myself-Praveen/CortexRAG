package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaProvider struct {
	url    string
	model  string
	client *http.Client
}

func NewOllamaProvider(url, model string) *OllamaProvider {
	if url == "" {
		url = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3"
	}
	return &OllamaProvider{
		url:    url,
		model:  model,
		client: &http.Client{},
	}
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (o *OllamaProvider) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody, _ := json.Marshal(ollamaGenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: false,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", o.url+"/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: %d - %s", resp.StatusCode, string(body))
	}

	var res ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Response, nil
}

func (o *OllamaProvider) GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	chunks := make(chan string, 50)
	errs := make(chan error, 1)

	go func() {
		defer close(chunks)
		defer close(errs)

		reqBody, _ := json.Marshal(ollamaGenerateRequest{
			Model:  o.model,
			Prompt: prompt,
			Stream: true,
		})

		req, err := http.NewRequestWithContext(ctx, "POST", o.url+"/api/generate", bytes.NewReader(reqBody))
		if err != nil {
			errs <- err
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := o.client.Do(req)
		if err != nil {
			errs <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errs <- fmt.Errorf("ollama API error: %d - %s", resp.StatusCode, string(body))
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var res ollamaGenerateResponse
			if err := json.Unmarshal(scanner.Bytes(), &res); err != nil {
				errs <- err
				return
			}
			chunks <- res.Response
			if res.Done {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			errs <- err
		}
	}()

	return chunks, errs
}
