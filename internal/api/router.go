package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", handler.Health)

	// API v1
	v1 := router.Group("/api/v1")
	{
		v1.POST("/search", handler.Search)
		v1.POST("/index-documents", handler.IndexDocuments)
	}

	return router
}
