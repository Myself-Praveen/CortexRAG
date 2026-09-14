package document

import "context"

// Extractor defines the interface for parsing documents.
type Extractor interface {
	// Extract parses a file and returns a Document.
	Extract(ctx context.Context, path string) (*Document, error)
}
