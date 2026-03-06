package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	EmbeddingProvider    string
	EmbeddingServiceURL  string
	OpenAIAPIKey         string
	OpenAIModel          string
	APIPort              string
	FAISSIndexPath       string
	FAISSDimension       int
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		EmbeddingProvider:   getEnv("EMBEDDING_PROVIDER", "local"),
		EmbeddingServiceURL: getEnv("EMBEDDING_SERVICE_URL", "http://localhost:8001"),
		OpenAIAPIKey:        getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:         getEnv("OPENAI_MODEL", "text-embedding-3-small"),
		APIPort:             getEnv("API_PORT", "8080"),
		FAISSIndexPath:      getEnv("FAISS_INDEX_PATH", "./data/faiss.index"),
		FAISSDimension:      384, // all-MiniLM-L6-v2 dimension
	}

	// Validate configuration
	if cfg.EmbeddingProvider == "openai" && cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required when using openai provider")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
