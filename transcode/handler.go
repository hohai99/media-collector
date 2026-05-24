package transcode

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler provides HTTP endpoints for video streaming via HLS.
type Handler struct {
	transcoder *Transcoder
}

// NewHandler creates a stream handler backed by the given Transcoder.
func NewHandler(t *Transcoder) *Handler {
	return &Handler{transcoder: t}
}

// ServeStream handles GET /stream?path=<encoded_path>
// It triggers a transcode (if needed) and redirects to the HLS playlist.
func (h *Handler) ServeStream(w http.ResponseWriter, r *http.Request) {
	videoPath := r.URL.Query().Get("path")
	if videoPath == "" {
		http.Error(w, "missing 'path' parameter", http.StatusBadRequest)
		return
	}

	// Validate file exists
	if _, err := os.Stat(videoPath); err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	job, err := h.transcoder.EnsureReady(r.Context(), videoPath)
	if err != nil {
		if r.Context().Err() != nil {
			return // client disconnected
		}
		http.Error(w, fmt.Sprintf("transcode failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect to the HLS playlist
	playlistURL := fmt.Sprintf("/stream/hls/%s/playlist.m3u8", job.ID)
	http.Redirect(w, r, playlistURL, http.StatusTemporaryRedirect)
}

// ServeHLS handles GET /stream/hls/<hash>/<file>
// It serves m3u8 playlist files and .ts segment files from the cache.
func (h *Handler) ServeHLS(w http.ResponseWriter, r *http.Request) {
	// Parse: /stream/hls/<hash>/<filename>
	path := strings.TrimPrefix(r.URL.Path, "/stream/hls/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	hash := parts[0]
	filename := parts[1]

	// Security: prevent path traversal
	if strings.Contains(hash, "..") || strings.Contains(filename, "..") ||
		strings.Contains(hash, "/") || strings.Contains(hash, "\\") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.transcoder.cfg.CacheDir, hash, filename)

	// Set appropriate content type
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".m3u8":
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	case ".ts":
		w.Header().Set("Content-Type", "video/MP2T")
	default:
		http.Error(w, "unsupported file type", http.StatusBadRequest)
		return
	}

	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeFile(w, r, filePath)
}
