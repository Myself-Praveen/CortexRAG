package llm

import (
	"fmt"
	"strings"

	"github.com/Myself-Praveen/CortexRAG/internal/vectorstore"
)

const systemPrompt = `You are a helpful assistant. Use the following context to answer the user's question. 
If the context doesn't contain the answer, say "I don't know based on the provided context."

Context:
%s

Question:
%s

Answer:`

func BuildPrompt(query string, results []vectorstore.SearchResult) string {
	var sb strings.Builder
	for i, res := range results {
		sb.WriteString(fmt.Sprintf("[%d] %s\n\n", i+1, res.Chunk.Content))
	}
	
	return fmt.Sprintf(systemPrompt, sb.String(), query)
}
