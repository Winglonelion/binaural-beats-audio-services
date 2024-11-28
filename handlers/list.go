package handlers

import (
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AudioFile struct {
	Name         string    `json:"name"`          // Name of the file
	Size         int64     `json:"size"`          // File size in bytes
	LastModified time.Time `json:"last_modified"` // Last modification time of the file
}

func ListFilesWithPagination(c *gin.Context) {
	dirPath := "./audio_files" // Directory containing audio files

	// Parse query parameters
	cursorParam := c.DefaultQuery("cursor", "") // Optional cursor for pagination
	limitParam := c.DefaultQuery("limit", "10") // Maximum number of files to return (default: 10)

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		c.JSON(400, gin.H{"error": "Invalid limit parameter"})
		return
	}

	var cursor time.Time
	if cursorParam != "" {
		cursor, err = time.Parse(time.RFC3339, cursorParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid cursor format"})
			return
		}
	}

	// Read files from the directory
	files := []AudioFile{}
	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || info.Name()[0] == '.' {
			return nil
		}

		// Append file metadata to the list
		files = append(files, AudioFile{
			Name:         info.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to read audio files"})
		return
	}

	// Sort files by LastModified in descending order (newest to oldest)
	sort.Slice(files, func(i, j int) bool {
		return files[i].LastModified.After(files[j].LastModified)
	})

	// Filter files based on the cursor
	filteredFiles := []AudioFile{}
	for _, file := range files {
		log.Printf("Checking file: %s, LastModified: %s, Cursor: %s", file.Name, file.LastModified, cursor)
		if cursorParam == "" || file.LastModified.Before(cursor) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Apply the limit to the filtered files
	paginatedFiles := filteredFiles
	if len(filteredFiles) > limit {
		paginatedFiles = filteredFiles[:limit]
	}

	// Determine the next cursor
	var nextCursor string
	if len(paginatedFiles) > 0 && len(filteredFiles) > limit {
		nextCursor = paginatedFiles[len(paginatedFiles)-1].LastModified.Format(time.RFC3339)
	}

	// Return the paginated list of files as JSON
	c.JSON(200, gin.H{
		"files":       paginatedFiles,
		"next_cursor": nextCursor,
	})
}
