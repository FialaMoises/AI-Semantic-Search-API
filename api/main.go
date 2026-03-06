package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yourusername/ai-semantic-search/internal/api"
	"github.com/yourusername/ai-semantic-search/internal/config"
	"github.com/yourusername/ai-semantic-search/internal/embeddings"
	"github.com/yourusername/ai-semantic-search/internal/search"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting AI Semantic Search API")
	log.Printf("Embedding Provider: %s", cfg.EmbeddingProvider)

	// Create embeddings provider
	embedder, err := embeddings.NewProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create embeddings provider: %v", err)
	}

	// Ensure data directory exists
	dataDir := filepath.Dir(cfg.FAISSIndexPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Create vector store
	vectorStore, err := search.NewVectorStore(cfg.FAISSDimension, cfg.FAISSIndexPath)
	if err != nil {
		log.Fatalf("Failed to create vector store: %v", err)
	}

	// Setup API
	handler := api.NewHandler(vectorStore, embedder)
	router := api.SetupRouter(handler)

	// Start server
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.APIPort)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
