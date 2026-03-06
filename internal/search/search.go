package search

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
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

// VectorStore manages the vector index and documents
type VectorStore struct {
	vectors   [][]float32
	documents []Document
	dimension int
	indexPath string
	mu        sync.RWMutex
}

// NewVectorStore creates a new vector store
func NewVectorStore(dimension int, indexPath string) (*VectorStore, error) {
	vs := &VectorStore{
		vectors:   make([][]float32, 0),
		documents: make([]Document, 0),
		dimension: dimension,
		indexPath: indexPath,
	}

	if err := vs.Load(); err != nil {
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

	for i, emb := range embeddings {
		if len(emb) != vs.dimension {
			return fmt.Errorf("embedding dimension mismatch: got %d, expected %d", len(emb), vs.dimension)
		}
		vs.vectors = append(vs.vectors, emb)
		vs.documents = append(vs.documents, docs[i])
		fmt.Printf("Added document %d: %s (ID: %s)\n", i, truncate(docs[i].Text, 50), docs[i].ID)
	}

	return nil
}

// Search performs similarity search using cosine similarity
func (vs *VectorStore) Search(queryEmbedding []float32, topK int) ([]SearchResult, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(queryEmbedding) != vs.dimension {
		return nil, fmt.Errorf("query embedding dimension mismatch: got %d, expected %d", len(queryEmbedding), vs.dimension)
	}

	if len(vs.vectors) == 0 {
		return []SearchResult{}, nil
	}

	type scoredResult struct {
		index int
		score float32
	}

	scores := make([]scoredResult, len(vs.vectors))
	for i, vec := range vs.vectors {
		scores[i] = scoredResult{
			index: i,
			score: cosineSimilarity(queryEmbedding, vec),
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	limit := topK
	if limit > len(scores) {
		limit = len(scores)
	}

	results := make([]SearchResult, limit)
	for i := 0; i < limit; i++ {
		results[i] = SearchResult{
			Document: vs.documents[scores[i].index],
			Score:    scores[i].score,
		}
	}

	return results, nil
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// Save persists the index and documents to disk
func (vs *VectorStore) Save() error {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	data, err := json.Marshal(map[string]interface{}{
		"vectors":   vs.vectors,
		"documents": vs.documents,
		"dimension": vs.dimension,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(vs.indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	fmt.Printf("Saved index with %d documents\n", len(vs.documents))
	return nil
}

// Load loads the index and documents from disk
func (vs *VectorStore) Load() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	data, err := os.ReadFile(vs.indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	var indexData struct {
		Vectors   [][]float32 `json:"vectors"`
		Documents []Document  `json:"documents"`
		Dimension int         `json:"dimension"`
	}

	if err := json.Unmarshal(data, &indexData); err != nil {
		return fmt.Errorf("failed to unmarshal index: %w", err)
	}

	vs.vectors = indexData.Vectors
	vs.documents = indexData.Documents
	vs.dimension = indexData.Dimension

	fmt.Printf("Loaded index with %d documents\n", len(vs.documents))
	return nil
}

// GetDocumentCount returns the number of indexed documents
func (vs *VectorStore) GetDocumentCount() int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return len(vs.documents)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}