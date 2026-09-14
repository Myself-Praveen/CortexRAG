package llm

import (
	"testing"

	"github.com/Myself-Praveen/CortexRAG/internal/document"
	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
	"github.com/stretchr/testify/assert"
)

func TestBuildPrompt(t *testing.T) {
	results := []vectorstore.SearchResult{
		{
			Chunk: &document.Chunk{Content: "Go is an open source programming language."},
		},
		{
			Chunk: &document.Chunk{Content: "It makes it easy to build simple, reliable, and efficient software."},
		},
	}

	prompt := BuildPrompt("What is Go?", results)

	assert.Contains(t, prompt, "Go is an open source")
	assert.Contains(t, prompt, "reliable, and efficient")
	assert.Contains(t, prompt, "What is Go?")
}
