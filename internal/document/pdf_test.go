package document

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPDFExtractor(t *testing.T) {
	extractor := NewPDFExtractor()

	t.Run("Extract non-existent file", func(t *testing.T) {
		_, err := extractor.Extract(context.Background(), "does-not-exist.pdf")
		assert.Error(t, err)
	})

	// To fully test this, we would need a dummy PDF file in testdata/
	// For now, we test the error handling.
}
