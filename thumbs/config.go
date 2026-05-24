package thumbs

import (
	"os"
	"path/filepath"
)

// Config holds all configurable settings for the thumbnail subsystem.
type Config struct {
	// MemoryMaxEntries is the maximum number of thumbnails held in the
	// in-process LRU cache (Tier 1). Each entry is a compressed JPEG []byte.
	MemoryMaxEntries int

	// DiskMaxBytes is the maximum total size of the on-disk thumbnail cache
	// (Tier 2). When exceeded, the least-recently-accessed entries are evicted.
	DiskMaxBytes int64

	// DiskCacheDir is the directory for the disk cache.
	// Defaults to {exeDir}/cache/thumbs if empty.
	DiskCacheDir string

	// DefaultWidth is the default thumbnail width in pixels.
	DefaultWidth int

	// DefaultHeight is the default thumbnail height in pixels.
	DefaultHeight int

	// DefaultFit controls default resize mode: "cover" or "contain".
	DefaultFit string

	// JPEGQuality is the JPEG encoding quality (1-100).
	JPEGQuality int

	// MaxWorkers is the maximum number of concurrent thumbnail generation
	// goroutines. Controlled via a buffered-channel semaphore.
	MaxWorkers int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MemoryMaxEntries: 500,
		DiskMaxBytes:     2 * 1024 * 1024 * 1024, // 2 GB
		DiskCacheDir:     defaultCacheDir(),
		DefaultWidth:     300,
		DefaultHeight:    300,
		DefaultFit:       "cover",
		JPEGQuality:      75,
		MaxWorkers:       4,
	}
}

// defaultCacheDir resolves {exeDir}/cache/thumbs.
func defaultCacheDir() string {
	exe, err := os.Executable()
	if err != nil {
		return filepath.Join(".", "cache", "thumbs")
	}
	return filepath.Join(filepath.Dir(exe), "cache", "thumbs")
}
