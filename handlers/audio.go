package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// StreamAudio streams audio files with Range Header support
func StreamAudio(context *gin.Context) {
	// Get the file name from the client (including extension)
	filename := context.Param("id")
	if filename == "" || !isValidFilename(filename) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}

	// Construct the full path to the file
	filePath := "./audio_files/" + filename

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Audio file not found"})
		return
	}
	defer file.Close()

	// Get file information
	fileStat, err := file.Stat()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve file info"})
		return
	}

	// Check Range Header (if present)
	rangeHeader := context.GetHeader("Range")
	if rangeHeader == "" {
		// No Range Header: Stream the entire file
		context.Writer.Header().Set("Content-Type", "audio/mpeg")
		context.Writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
		http.ServeContent(context.Writer, context.Request, fileStat.Name(), fileStat.ModTime(), file)
		return
	}

	// Parse the Range Header
	byteRange, err := parseRangeHeader(rangeHeader, fileStat.Size())
	if err != nil {
		context.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": err.Error()})
		return
	}

	// Configure headers for partial content response
	context.Writer.Header().Set("Content-Type", "audio/mpeg")
	context.Writer.Header().Set("Content-Range", "bytes "+strconv.FormatInt(byteRange.Start, 10)+"-"+strconv.FormatInt(byteRange.End, 10)+"/"+strconv.FormatInt(fileStat.Size(), 10))
	context.Writer.Header().Set("Accept-Ranges", "bytes")
	context.Writer.Header().Set("Content-Length", strconv.FormatInt(byteRange.End-byteRange.Start+1, 10))
	context.Writer.WriteHeader(http.StatusPartialContent)

	// Stream the requested range
	file.Seek(byteRange.Start, 0)
	buffer := make([]byte, 1024*16) // Buffer size: 16KB
	bytesToRead := byteRange.End - byteRange.Start + 1

	for bytesToRead > 0 {
		if bytesToRead < int64(len(buffer)) {
			buffer = buffer[:bytesToRead]
		}
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
			return
		}
		if n == 0 {
			break
		}
		context.Writer.Write(buffer[:n])
		bytesToRead -= int64(n)
	}
}

// isValidFilename validates the filename provided by the client
func isValidFilename(filename string) bool {
	// Ensure the file name does not contain invalid characters and has an extension
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return false
	}

	// Check allowed file extensions (e.g., only accept .mp3, .wav, ...)
	allowedExtensions := []string{".mp3", ".wav", ".flac"}
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

// parseRangeHeader parses the Range Header and returns the byte range
func parseRangeHeader(rangeHeader string, fileSize int64) (struct{ Start, End int64 }, error) {
	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		return struct{ Start, End int64 }{}, errors.New("invalid range header")
	}

	byteRange := strings.Split(parts[1], "-")
	start, err := strconv.ParseInt(byteRange[0], 10, 64)
	if err != nil || start < 0 {
		return struct{ Start, End int64 }{}, errors.New("invalid start range")
	}

	end := fileSize - 1
	if len(byteRange) > 1 && byteRange[1] != "" {
		end, err = strconv.ParseInt(byteRange[1], 10, 64)
		if err != nil || end >= fileSize || end < start {
			return struct{ Start, End int64 }{}, errors.New("invalid end range")
		}
	}

	return struct{ Start, End int64 }{Start: start, End: end}, nil
}
