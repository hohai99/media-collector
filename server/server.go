package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"media-collector/thumbs"
	"media-collector/transcode"
)

// FileServer is a standalone HTTP server that serves local files,
// thumbnails, and video streams outside the Wails asset pipeline.
type FileServer struct {
	port       int
	listener   net.Listener
	httpServer *http.Server
	thumbGen   *thumbs.Generator
	transcoder *transcode.Transcoder
	metrics    *Metrics
	logger     *slog.Logger
}

// New creates a FileServer. Call Start() to begin accepting connections.
// transcoder may be nil if ffmpeg is not available.
func New(thumbGen *thumbs.Generator, tc *transcode.Transcoder, logger *slog.Logger) (*FileServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	s := &FileServer{
		port:       listener.Addr().(*net.TCPAddr).Port,
		listener:   listener,
		thumbGen:   thumbGen,
		transcoder: tc,
		metrics:    NewMetrics(),
		logger:     logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/file/", s.handleFile)
	mux.HandleFunc("/thumb", s.handleThumb)
	mux.Handle("/metrics", s.metrics)

	// Video transcode routes (only if transcoder is available)
	if tc != nil {
		handler := transcode.NewHandler(tc)
		mux.HandleFunc("/stream", handler.ServeStream)
		mux.HandleFunc("/stream/hls/", handler.ServeHLS)
	}

	s.httpServer = &http.Server{
		Handler: s.withMiddleware(mux),
	}

	return s, nil
}

// Start begins serving in a background goroutine.
func (s *FileServer) Start() {
	s.logger.Info("file server starting", "port", s.port)
	go func() {
		if err := s.httpServer.Serve(s.listener); err != nil && err != http.ErrServerClosed {
			s.logger.Error("file server error", "error", err)
		}
	}()
}

// Shutdown gracefully stops the server with a 5-second timeout,
// and kills any running transcode processes.
func (s *FileServer) Shutdown(ctx context.Context) error {
	if s.transcoder != nil {
		s.transcoder.Shutdown()
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s.logger.Info("file server shutting down")
	return s.httpServer.Shutdown(shutdownCtx)
}

// Port returns the port the server is listening on.
func (s *FileServer) Port() int {
	return s.port
}

// BaseURL returns the base URL for file serving (e.g. "http://127.0.0.1:12345/file/").
func (s *FileServer) BaseURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d/file/", s.port)
}

// StreamBaseURL returns the base URL for stream requests.
func (s *FileServer) StreamBaseURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d/stream", s.port)
}

// HasTranscoder returns true if video transcoding is available.
func (s *FileServer) HasTranscoder() bool {
	return s.transcoder != nil
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func (s *FileServer) handleFile(w http.ResponseWriter, r *http.Request) {
	done := s.metrics.RecordStart()
	defer done(false)

	// Strip the /file/ prefix
	rawPath := strings.TrimPrefix(r.URL.Path, "/file/")

	// URL-decode
	filePath, err := url.PathUnescape(rawPath)
	if err != nil {
		filePath = rawPath
	}

	// Convert to OS path
	filePath = filepath.FromSlash(filePath)

	// Open and stream the file
	f, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ServeContent handles Range requests, Content-Type detection, and
	// streams the file without loading it all into memory.
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
}

func (s *FileServer) handleThumb(w http.ResponseWriter, r *http.Request) {
	done := s.metrics.RecordStart()

	start := time.Now()
	s.thumbGen.ServeHTTP(w, r)

	cacheHit := time.Since(start) < 5*time.Millisecond
	done(cacheHit)
}

// ---------------------------------------------------------------------------
// Middleware
// ---------------------------------------------------------------------------

type middlewareHandler struct {
	mux    *http.ServeMux
	logger *slog.Logger
}

func (s *FileServer) withMiddleware(mux *http.ServeMux) http.Handler {
	return &middlewareHandler{mux: mux, logger: s.logger}
}

func (h *middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Range")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Range, Content-Length")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	start := time.Now()
	h.mux.ServeHTTP(w, r)

	// Structured request log (skip /metrics to avoid noise)
	if r.URL.Path != "/metrics" {
		h.logger.Debug("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}
