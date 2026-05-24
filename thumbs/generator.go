package thumbs

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"golang.org/x/sync/singleflight"
)

// Generator handles thumbnail creation with request deduplication and
// concurrency limiting.
type Generator struct {
	cache  *ThumbnailCache
	sem    chan struct{}
	group  singleflight.Group
	cfg    Config
	logger *slog.Logger
}

// NewGenerator creates a Generator backed by the given cache.
func NewGenerator(cfg Config, logger *slog.Logger) (*Generator, error) {
	cache, err := NewThumbnailCache(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("init thumbnail cache: %w", err)
	}

	return &Generator{
		cache:  cache,
		sem:    make(chan struct{}, cfg.MaxWorkers),
		cfg:    cfg,
		logger: logger,
	}, nil
}

// Generate returns a JPEG thumbnail for the image at the given path.
// It uses a two-tier cache, singleflight deduplication, and a worker-pool
// semaphore to bound concurrency.
func (g *Generator) Generate(ctx context.Context, absPath string, w, h int, fit string) ([]byte, error) {
	start := time.Now()

	// Stat the source file for mod-time (used in cache key).
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("stat source: %w", err)
	}

	cacheKey := CacheKey(absPath, w, h, fit, info.ModTime())
	diskKey := DiskKey(absPath, w, h, fit)

	// Fast path: cache hit
	if data, ok := g.cache.Get(cacheKey, diskKey); ok {
		g.logger.Debug("thumbnail cache hit",
			"path", absPath,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return data, nil
	}

	// Slow path: generate via singleflight (dedup concurrent requests).
	result, err, _ := g.group.Do(cacheKey, func() (interface{}, error) {
		return g.generate(ctx, absPath, w, h, fit, info.ModTime(), cacheKey, diskKey)
	})
	if err != nil {
		return nil, err
	}

	data := result.([]byte)
	g.logger.Debug("thumbnail generated",
		"path", absPath,
		"w", w, "h", h, "fit", fit,
		"size", len(data),
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return data, nil
}

// generate does the actual decode → resize → encode work.
func (g *Generator) generate(
	ctx context.Context,
	absPath string,
	w, h int,
	fit string,
	modTime time.Time,
	cacheKey, diskKey string,
) ([]byte, error) {
	// Acquire semaphore slot (bounded concurrency).
	select {
	case g.sem <- struct{}{}:
		defer func() { <-g.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Check cancellation before heavy work.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Decode source image.
	src, err := imaging.Open(absPath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decode image %q: %w", absPath, err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Resize.
	var resized *image.NRGBA
	switch strings.ToLower(fit) {
	case "contain":
		resized = imaging.Fit(src, w, h, imaging.Lanczos)
	default: // "cover"
		resized = imaging.Fill(src, w, h, imaging.Center, imaging.Lanczos)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Encode to JPEG.
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: g.cfg.JPEGQuality}); err != nil {
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}

	data := buf.Bytes()

	// Store in cache (fire-and-forget for disk errors).
	if err := g.cache.Put(cacheKey, diskKey, data, modTime); err != nil {
		g.logger.Warn("failed to cache thumbnail", "path", absPath, "error", err)
	}

	return data, nil
}

// Cache returns the underlying ThumbnailCache (for metrics).
func (g *Generator) Cache() *ThumbnailCache {
	return g.cache
}

// isImagePath checks whether the file extension is a supported image type.
func isImagePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tif", ".tiff":
		return true
	}
	return false
}
