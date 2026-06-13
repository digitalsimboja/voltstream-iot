package main

import (
	"log/slog"
	"os"

	"github.com/digitalsimboja/voltstream-iot/ingestion-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "voltstream-ingestion-api"})
	})

	r.POST("/telemetry", handlers.IngestTelemetry)

	// TODO: register WebSocket endpoint for dashboard live streaming
	// r.GET("/ws", handlers.StreamWebSocket)

	slog.Info("ingestion API starting", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
