package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ai-semantic-search/internal/embeddings"
	"github.com/yourusername/ai-semantic-search/internal/search"
)

type Handler struct {
	vectorStore *search.VectorStore
	embedder    embeddings.Provider
}

func NewHandler(vectorStore *search.VectorStore, embedder embeddings.Provider) *Handler {
	return &Handler{
		vectorStore: vectorStore,
		embedder:    embedder,
	}
}

// SearchRequest represents a search query
type SearchRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

// IndexRequest represents documents to be indexed
type IndexRequest struct {
	Documents []search.Document `json:"documents" binding:"required"`
}

// Search handles semantic search requests
func (h *Handler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TopK == 0 {
		req.TopK = 5
	}

	embedding, err := h.embedder.GenerateEmbedding(c.Request.Context(), req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate embedding: %v", err)})
		return
	}

	results, err := h.vectorStore.Search(embedding, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("search failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   req.Query,
		"results": results,
		"count":   len(results),
	})
}

// IndexDocuments handles document indexing requests
func (h *Handler) IndexDocuments(c *gin.Context) {
	var req IndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Documents) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no documents provided"})
		return
	}

	texts := make([]string, len(req.Documents))
	for i, doc := range req.Documents {
		texts[i] = doc.Text
	}

	embeddings, err := h.embedder.GenerateBatchEmbeddings(c.Request.Context(), texts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate embeddings: %v", err)})
		return
	}

	if err := h.vectorStore.AddDocuments(req.Documents, embeddings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to index documents: %v", err)})
		return
	}

	if err := h.vectorStore.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save index: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "documents indexed successfully",
		"indexed_count":  len(req.Documents),
		"total_documents": h.vectorStore.GetDocumentCount(),
	})
}

// Health checks the API health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":          "healthy",
		"document_count":  h.vectorStore.GetDocumentCount(),
		"embedding_dimension": h.embedder.GetDimension(),
	})
}
