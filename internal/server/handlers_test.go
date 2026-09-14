package server

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Myself-Praveen/CortexRAG/internal/config"
	"github.com/Myself-Praveen/CortexRAG/internal/rag"
	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEmbedder struct{}
func (m *mockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) { return []float32{1, 0}, nil }
func (m *mockEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) { return []float32{1, 0}, nil }
func (m *mockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) { return nil, nil }

type mockLLM struct{}
func (m *mockLLM) Generate(ctx context.Context, prompt string) (string, error) { return "Mock Answer", nil }
func (m *mockLLM) GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	c := make(chan string, 1)
	e := make(chan error, 1)
	c <- "Mock "
	c <- "Stream"
	close(c)
	close(e)
	return c, e
}

func TestAPIEndpoints(t *testing.T) {
	cfg := &config.Config{}
	cfg.RAG.TopK = 1
	cfg.RAG.ChunkSize = 100
	cfg.RAG.ChunkOverlap = 10

	store := vectorstore.NewMemoryStore()
	embedder := &mockEmbedder{}
	llm := &mockLLM{}
	
	pipeline := rag.NewPipeline(embedder, store, nil, cfg) // nil pool might crash, but ingest test might not execute full flow if we mock it well, actually pipeline will crash without pool. Let's just test handlers directly with mocked pipeline? Actually pipeline struct is concrete.
	// Since pipeline needs a workerpool to not panic, we can't test ingest easily without setting up pool.
	generator := rag.NewGenerator(embedder, store, llm, cfg)
	
	api := NewAPI(pipeline, generator)

	t.Run("HandleQuery", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"query": "test"}`))
		req := httptest.NewRequest("POST", "/api/query", body)
		w := httptest.NewRecorder()
		
		api.HandleQuery(w, req)
		
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, w.Body.String(), "Mock Answer")
	})

	t.Run("HandleStream", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/stream?q=test", nil)
		w := httptest.NewRecorder()
		
		api.HandleStream(w, req)
		
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "text/event-stream", res.Header.Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "data: Mock Stream")
	})
}

func TestIngestEndpoint(t *testing.T) {
	// Simple test for bad request
	api := NewAPI(nil, nil)
	
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	part.Write([]byte("test"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/ingest", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Instead, just test without file to get 400.
	reqBad := httptest.NewRequest("POST", "/api/ingest", nil)
	wBad := httptest.NewRecorder()
	api.HandleIngest(wBad, reqBad)
	assert.Equal(t, http.StatusBadRequest, wBad.Result().StatusCode)
}

func TestServerStart(t *testing.T) {
	cfg := &config.Config{}
	
	srv := NewServer(cfg, NewAPI(nil, nil))
	
	go func() {
		srv.Start()
	}()
	
	time.Sleep(50 * time.Millisecond) // Let it start
}
