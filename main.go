package main

import (
	"binaural_beats_audio_services/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Endpoint to list all audio files
	router.GET("/api/audio", handlers.ListFilesWithPagination)

	// Endpoint to stream audio
	router.GET("/api/audio/:id", handlers.StreamAudio)

	router.Run(":8080") // Run server on port 8080
}
