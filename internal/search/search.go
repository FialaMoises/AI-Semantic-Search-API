package search

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/DataIntelligenceCrew/go-faiss"
)

// Document represents a searchable document
type Document struct {
	ID       string                 `json:"id"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SearchResult represents a search result with score
type SearchResult struct {
	Document Document `json:"document"`
	Score    float32  `json:"score"`
}

// VectorStore manages the FAISS index and documents
type VectorStore struct {
	index      *faiss.IndexImpl
	documents  map[int64]Document
	dimension  int
	indexPath  string
	mu         sync.RWMutex
	nextID     int64
}

// NewVectorStore creates a new vector store
func NewVectorStore(dimension int, indexPath string) (*VectorStore, error) {
	// Create index (L2 distance)
	index, err := faiss.NewIndexFlatL2(dimension)
	if err != nil {
		return nil, fmt.Errorf("failed to create FAISS index: %w", err)
	}

	vs := &VectorStore{
		index:     index,
		documents: make(map[int64]Document),
		dimension: dimension,
		indexPath: indexPath,
		nextID:    0,
	}

	// Try to load existing index
	if err := vs.Load(); err != nil {
		// If load fails, just use empty index (not an error for new deployments)
		fmt.Printf("Starting with fresh index: %v\n", err)
	}

	return vs, nil
}

// AddDocuments adds documents with their embeddings to the index
func (vs *VectorStore) AddDocuments(docs []Document, embeddings [][]float32) error {
	if len(docs) != len(embeddings) {
		return fmt.Errorf("number of documents (%d) must match number of embeddings (%d)", len(docs), len(embeddings))
	}

	vs.mu.Lock()
	defer vs.mu.Unlock()

	// Prepare vectors for FAISS
	vectors := make([]float32, 0, len(embeddings)*vs.dimension)
	for _, emb := range embeddings {
		if len(emb) != vs.dimension {
			return fmt.Errorf("embedding dimension mismatch: got %d, expected %d", len(emb), vs.dimension)
		}
		vectors = append(vectors, emb...)
	}

	// Add to index
	if err := vs.index.Add(vectors); err != nil {
		return fmt.Errorf("failed to add to FAISS index: %w", err)
	}

	// Store documents
	for i, doc := range docs {
		vs.documents[vs.nextID] = doc
		vs.nextID++
		fmt.Printf("Added document %d: %s (ID: %s)\n", i, doc.Text[:min(50, len(doc.Text))], doc.ID)
	}

	return nil
}

// Search performs similarity search
func (vs *VectorStore) Search(queryEmbedding []float32, topK int) ([]SearchResult, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(queryEmbedding) != vs.dimension {
		return nil, fmt.Errorf("query embedding dimension mismatch: got %d, expected %d", len(queryEmbedding), vs.dimension)
	}

	// Perform search
	distances, ids, err := vs.index.Search(queryEmbedding, int64(topK))
	if err != nil {
		return nil, fmt.Errorf("FAISS search failed: %w", err)
	}

	// Build results
	results := make([]SearchResult, 0, len(ids))
	for i, id := range ids {
		if id == -1 {
			// FAISS returns -1 for empty slots
			continue
		}

		doc, exists := vs.documents[id]
		if !exists {
			continue
		}

		results = append(results, SearchResult{
			Document: doc,
			Score:    distances[i],
		})
	}

	return results, nil
}

// Save persists the index and documents to disk
func (vs *VectorStore) Save() error {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	// Save FAISS index
	if err := faiss.WriteIndex(vs.index, vs.indexPath); err != nil {
		return fmt.Errorf("failed to save FAISS index: %w", err)
	}

	// Save documents metadata
	metadataPath := vs.indexPath + ".meta"
	data, err := json.Marshal(map[string]interface{}{
		"documents": vs.documents,
		"nextID":    vs.nextID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	fmt.Printf("Saved index with %d documents\n", len(vs.documents))
	return nil
}

// Load loads the index and documents from disk
func (vs *VectorStore) Load() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	// Load FAISS index
	index, err := faiss.ReadIndex(vs.indexPath)
	if err != nil {
		return fmt.Errorf("failed to load FAISS index: %w", err)
	}
	vs.index = index

	// Load documents metadata
	metadataPath := vs.indexPath + ".meta"
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata struct {
		Documents map[string]Document `json:"documents"`
		NextID    int64               `json:"nextID"`
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Convert string keys back to int64
	vs.documents = make(map[int64]Document)
	for k, v := range metadata.Documents {
		var id int64
		fmt.Sscanf(k, "%d", &id)
		vs.documents[id] = v
	}
	vs.nextID = metadata.NextID

	fmt.Printf("Loaded index with %d documents\n", len(vs.documents))
	return nil
}

// GetDocumentCount returns the number of indexed documents
func (vs *VectorStore) GetDocumentCount() int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return len(vs.documents)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
