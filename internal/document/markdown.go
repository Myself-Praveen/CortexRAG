package document

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// MarkdownExtractor implements the Extractor interface for .md files.
type MarkdownExtractor struct{}

func NewMarkdownExtractor() *MarkdownExtractor {
	return &MarkdownExtractor{}
}

func (m *MarkdownExtractor) Extract(ctx context.Context, path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	content := string(data)
	// Simple cleanup: standardizing newlines
	content = strings.ReplaceAll(content, "\r\n", "\n")

	// Get base name as source
	source := filepath.Base(path)

	doc := &Document{
		ID:      uuid.NewString(),
		Source:  source,
		Content: content,
		Metadata: Metadata{
			"type": "markdown",
			"path": path,
		},
	}

	return doc, nil
}
