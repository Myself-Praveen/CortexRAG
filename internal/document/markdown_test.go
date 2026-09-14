package document

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownExtractor(t *testing.T) {
	extractor := NewMarkdownExtractor()

	t.Run("Extract valid markdown file", func(t *testing.T) {
		// Setup temporary markdown file
		dir := t.TempDir()
		path := filepath.Join(dir, "test.md")
		content := "# Header\n\nSome content."
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		// Act
		doc, err := extractor.Extract(context.Background(), path)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, doc.ID)
		assert.Equal(t, "test.md", doc.Source)
		assert.Equal(t, content, doc.Content)
		assert.Equal(t, "markdown", doc.Metadata["type"])
		assert.Equal(t, path, doc.Metadata["path"])
	})

	t.Run("Extract non-existent file", func(t *testing.T) {
		_, err := extractor.Extract(context.Background(), "does-not-exist.md")
		assert.Error(t, err)
	})
}
