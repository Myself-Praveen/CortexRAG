package document

import (
	"strings"

	"github.com/google/uuid"
)

// Chunker is responsible for splitting text into smaller chunks.
type Chunker struct{}

func NewChunker() *Chunker {
	return &Chunker{}
}

// Chunk splits a Document into multiple Chunks based on size and overlap.
func (c *Chunker) Chunk(doc *Document, size, overlap int) []*Chunk {
	if size <= 0 {
		size = 512
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= size {
		overlap = size / 2
	}

	content := doc.Content
	var chunks []*Chunk
	var index int

	// Simple recursive character splitting strategy:
	// Split by paragraphs first, then sentences if paragraphs are too long.
	// For simplicity in this implementation, we will use a basic sliding window
	// of characters/words. A more advanced version would use tokenization.
	
	words := strings.Fields(content)
	if len(words) == 0 {
		return chunks
	}

	for i := 0; i < len(words); {
		end := i + size
		if end > len(words) {
			end = len(words)
		}

		chunkWords := words[i:end]
		chunkContent := strings.Join(chunkWords, " ")

		chunks = append(chunks, &Chunk{
			ID:      uuid.NewString(),
			DocID:   doc.ID,
			Content: chunkContent,
			Index:   index,
		})
		index++

		if end == len(words) {
			break
		}

		// Advance by size - overlap
		step := size - overlap
		if step <= 0 {
			step = 1
		}
		i += step
	}

	return chunks
}
