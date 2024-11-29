
# Binaural Beats Audio Backend

## Overview

This is the backend service for the Binaural Beats Audio App. It provides API endpoints for serving audio files, retrieving metadata, and supporting analytics. The backend is implemented using **Go** for high performance and scalability.

## Features

- Streaming audio files via API.
- Download audio files via API.
- Provide audio metadata, including:
  - Name of the audio file.
  - Author of the audio file.
  - Precomputed FFT data (for visualizations).
- Support error handling and logging.
- Easy scalability for high traffic.

## Requirements

- **Go** (>=1.20)
- Other dependencies managed via `go.mod`

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/Winglonelion/binaural-beats-audio-services.git
cd binaural-beats-audio-services
```

### 2. Quick start server with Docker
  ```bash
  docker-compose up -d
  ```

### 3. Install dependencies

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

- **Parameters**:
  - `cursor` (string): Cursor to pagination data, in ios string time format.
  - `limit` (number): Limit rows per page.

**Response:**
```json
{
  "files": [
    {
      "name": "audio_1.mp3",
      "size": 12345678,
      "last_modified": "2024-11-25T03:47:54Z",
      metadata: {
        id: "1",
        name: "Audio 1",
        author: "Nemo",
        fft: [],
        cover_image: "https://example.com/img1.png",
        thumbhash: "RanDomHazh"
      }
    },
    {
      "name": "audio_2.mp3",
      "size": 9876543,
      "last_modified": "2024-11-27T10:55:18Z",
      metadata: {
        id: "2",
        name: "Audio 2",
        author: "Fin",
        fft: [],
        cover_image: "https://example.com/img2.png",
        thumbhash: "RanDomHazh"
      }
    }
  ]
}
```


### 2. `GET /api/audio/:filename`
Stream an audio file.
- **Parameters**:
  - `filename` (string): Name of the audio file.

### 3. `GET /api/download/:filename`
Download audio


- **Parameters**:
  - `fileName` (string): Name of the audio file.


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
