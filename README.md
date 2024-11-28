
# Binaural Beats Audio Backend

## Overview

This is the backend service for the Binaural Beats Audio App. It provides API endpoints for serving audio files, retrieving metadata, and supporting analytics. The backend is implemented using **Go** for high performance and scalability.

## Features

- Serve audio files via API.
- Provide audio metadata, including:
  - Name of the audio file.
  - Author of the audio file.
  - Precomputed FFT data (for visualizations).
- Support error handling and logging.
- Easy scalability for high traffic.
- Extensible analytics integration.

## Requirements

- **Go** (>=1.20)
- **FFmpeg** (for processing audio, if needed)
- Other dependencies managed via `go.mod`

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/Winglonelion/binaural-beats-audio-services.git
cd binaural-beats-audio-services
```

### 2. Install dependencies

Dependencies are managed via `go.mod`. Run the following command to ensure all dependencies are installed:

```bash
go mod tidy
```

### 4. Run the server

Start the backend server:

```bash
go run main.go
```

The server will be available at `http://localhost:8080`.

## API Endpoints

### 1. `GET /api/audio`
Retrieve a list of available audio files.

**Response:**
```json
{
  "files": [
    {
      "name": "audio_1.mp3",
      "size": 12345678,
      "last_modified": "2024-11-25T03:47:54Z"
    },
    {
      "name": "audio_2.mp3",
      "size": 9876543,
      "last_modified": "2024-11-27T10:55:18Z"
    }
  ]
}
```

### 2. `GET /api/audio/metadata/:name`
Retrieve metadata for a specific audio file.

**Response:**
```json
{
  "name": "audio_1.mp3",
  "author": "Unknown Artist",
  "fft": [0.1, 0.2, 0.3, ...]
}
```

### 3. `GET /api/audio/:name`
Stream an audio file.

- **Parameters**:
  - `name` (string): Name of the audio file.


## Development

### 1. Run in Development Mode
To run the server with hot-reloading:

```bash
air
```

Install `air` if you haven't already:

```bash
go install github.com/cosmtrek/air@latest
```

### 2. Testing

Run unit tests:

```bash
go test ./...
```

### 3. Linting

Ensure code quality:

```bash
golangci-lint run
```

Install `golangci-lint` if not already installed:

```bash
brew install golangci-lint
```

## Deployment

To build the production binary:

```bash
go build -o binaural_beats_backend main.go
```

Then deploy the binary to your server.
