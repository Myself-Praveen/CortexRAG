package document

// Metadata represents standard document metadata.
type Metadata map[string]any

// Document represents a fully ingested document.
type Document struct {
	ID       string   `json:"id"`
	Source   string   `json:"source"`
	Content  string   `json:"content"`
	Metadata Metadata `json:"metadata"`
}

// Chunk represents a segment of a document.
type Chunk struct {
	ID        string    `json:"id"`
	DocID     string    `json:"doc_id"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding,omitempty"`
	Index     int       `json:"index"`
}
