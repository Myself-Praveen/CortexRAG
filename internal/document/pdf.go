package document

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
)

// PDFExtractor implements the Extractor interface for .pdf files.
type PDFExtractor struct{}

func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

func (p *PDFExtractor) Extract(ctx context.Context, path string) (*Document, error) {
	content, err := readPdf(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF %s: %w", path, err)
	}

	source := filepath.Base(path)

	doc := &Document{
		ID:      uuid.NewString(),
		Source:  source,
		Content: content,
		Metadata: Metadata{
			"type": "pdf",
			"path": path,
		},
	}

	return doc, nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var sb strings.Builder
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		
		text, err := p.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
