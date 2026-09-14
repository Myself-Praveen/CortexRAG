package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunker(t *testing.T) {
	chunker := NewChunker()

	doc := &Document{
		ID:      "doc-1",
		Content: "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10",
	}

	t.Run("Chunk with overlap", func(t *testing.T) {
		chunks := chunker.Chunk(doc, 5, 2)
		
		assert.Len(t, chunks, 3)
		assert.Equal(t, "word1 word2 word3 word4 word5", chunks[0].Content)
		assert.Equal(t, "word4 word5 word6 word7 word8", chunks[1].Content)
		assert.Equal(t, "word7 word8 word9 word10", chunks[2].Content)
		
		assert.Equal(t, "doc-1", chunks[0].DocID)
		assert.NotEmpty(t, chunks[0].ID)
		assert.Equal(t, 0, chunks[0].Index)
		assert.Equal(t, 1, chunks[1].Index)
	})

	t.Run("Chunk without overlap", func(t *testing.T) {
		chunks := chunker.Chunk(doc, 5, 0)
		
		assert.Len(t, chunks, 2)
		assert.Equal(t, "word1 word2 word3 word4 word5", chunks[0].Content)
		assert.Equal(t, "word6 word7 word8 word9 word10", chunks[1].Content)
	})

	t.Run("Chunk size greater than content", func(t *testing.T) {
		chunks := chunker.Chunk(doc, 20, 5)
		
		assert.Len(t, chunks, 1)
		assert.Equal(t, doc.Content, chunks[0].Content)
	})
}
