package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DownloadAudio handles the download of audio files
func DownloadAudio(c *gin.Context) {
	// Directory containing audio files
	audioDir := "./audio_files"

	// Get filename from the URL parameter
	fileName := c.Param("filename")

	// Construct the full file path
	filePath := filepath.Join(audioDir, fileName)

	// Check if the file exists and is accessible
	if !fileExists(filePath) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file as an attachment
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

// fileExists checks if a file exists and is not a directory
// fileExists checks if a file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path) // Use os.Stat instead of filepath.Stat
	if err != nil {
		return false
	}
	return !info.IsDir()
}
