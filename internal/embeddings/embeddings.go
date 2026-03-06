package embeddings

import (
	"context"
	"fmt"

	"github.com/yourusername/ai-semantic-search/internal/config"
)

// Provider defines the interface for embedding generation
type Provider interface {
	// GenerateEmbedding creates a vector embedding from text
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)

	// GenerateBatchEmbeddings creates embeddings for multiple texts
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)

	// GetDimension returns the dimension of the embeddings
	GetDimension() int
}

// NewProvider creates an embedding provider based on configuration
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.EmbeddingProvider {
	case "local":
		return NewLocalProvider(cfg.EmbeddingServiceURL)
	case "openai":
		return NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.OpenAIModel)
	default:
		return nil, fmt.Errorf("unknown embedding provider: %s", cfg.EmbeddingProvider)
	}
}
