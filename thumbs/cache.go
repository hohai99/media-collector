package thumbs

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// diskMeta is the sidecar metadata stored alongside each cached thumbnail.
type diskMeta struct {
	SrcModTime int64 `json:"srcModTime"` // source file mod-time (unix)
	AccessTime int64 `json:"accessTime"` // last access time (unix)
	Size       int64 `json:"size"`       // thumbnail file size in bytes
}

// ThumbnailCache implements a two-tier cache (memory LRU + disk).
type ThumbnailCache struct {
	mem     *lru.Cache[string, []byte]
	diskDir string
	maxDisk int64
	mu      sync.Mutex // protects disk eviction
	logger  *slog.Logger
}

// NewThumbnailCache creates and returns a ready-to-use ThumbnailCache.
func NewThumbnailCache(cfg Config, logger *slog.Logger) (*ThumbnailCache, error) {
	if err := os.MkdirAll(cfg.DiskCacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}

	mem, err := lru.New[string, []byte](cfg.MemoryMaxEntries)
	if err != nil {
		return nil, fmt.Errorf("create memory cache: %w", err)
	}

	return &ThumbnailCache{
		mem:     mem,
		diskDir: cfg.DiskCacheDir,
		maxDisk: cfg.DiskMaxBytes,
		logger:  logger,
	}, nil
}

// CacheKey builds a deterministic cache key from the image path, dimensions,
// fit mode, and source file modification time.
func CacheKey(absPath string, w, h int, fit string, modTime time.Time) string {
	raw := fmt.Sprintf("%s:%d:%d:%s:%d", absPath, w, h, fit, modTime.Unix())
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

// DiskKey is like CacheKey but omits modTime — used for the on-disk filename
// so that a changed file naturally produces a different CacheKey while the
// disk path stays stable per (path, size, fit) combination.
func DiskKey(absPath string, w, h int, fit string) string {
	raw := fmt.Sprintf("%s:%d:%d:%s", absPath, w, h, fit)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

// Get looks up a thumbnail. It checks memory first, then disk.
// Returns the JPEG bytes and true on hit, nil and false on miss.
func (c *ThumbnailCache) Get(cacheKey, diskKey string) ([]byte, bool) {
	// Tier 1 — memory
	if data, ok := c.mem.Get(cacheKey); ok {
		return data, true
	}

	// Tier 2 — disk
	imgPath := filepath.Join(c.diskDir, diskKey+".jpg")
	metaPath := filepath.Join(c.diskDir, diskKey+".meta")

	data, err := os.ReadFile(imgPath)
	if err != nil {
		return nil, false
	}

	// Update access time in metadata
	c.touchMeta(metaPath, int64(len(data)))

	// Promote to memory
	c.mem.Add(cacheKey, data)

	return data, true
}

// Put stores a thumbnail in both memory and disk caches.
func (c *ThumbnailCache) Put(cacheKey, diskKey string, data []byte, srcModTime time.Time) error {
	// Memory
	c.mem.Add(cacheKey, data)

	// Disk
	imgPath := filepath.Join(c.diskDir, diskKey+".jpg")
	metaPath := filepath.Join(c.diskDir, diskKey+".meta")

	if err := os.WriteFile(imgPath, data, 0o644); err != nil {
		return fmt.Errorf("write thumbnail: %w", err)
	}

	meta := diskMeta{
		SrcModTime: srcModTime.Unix(),
		AccessTime: time.Now().Unix(),
		Size:       int64(len(data)),
	}
	metaBytes, _ := json.Marshal(meta)
	if err := os.WriteFile(metaPath, metaBytes, 0o644); err != nil {
		// Non-fatal — the image is still cached.
		c.logger.Warn("failed to write cache metadata", "path", metaPath, "error", err)
	}

	// Async eviction check
	go c.evictIfNeeded()

	return nil
}

// EvictStale removes disk entries until total size is under maxDisk.
func (c *ThumbnailCache) evictIfNeeded() {
	c.mu.Lock()
	defer c.mu.Unlock()

	type entry struct {
		diskKey    string
		accessTime int64
		size       int64
	}

	var entries []entry
	var totalSize int64

	// Scan all .meta files
	filepath.WalkDir(c.diskDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".meta" {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var m diskMeta
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil
		}

		base := path[:len(path)-len(".meta")]
		dk := filepath.Base(base)
		entries = append(entries, entry{diskKey: dk, accessTime: m.AccessTime, size: m.Size})
		totalSize += m.Size
		return nil
	})

	if totalSize <= c.maxDisk {
		return
	}

	// Sort by access time ascending (oldest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].accessTime < entries[j].accessTime
	})

	for _, e := range entries {
		if totalSize <= c.maxDisk {
			break
		}
		imgPath := filepath.Join(c.diskDir, e.diskKey+".jpg")
		metaPath := filepath.Join(c.diskDir, e.diskKey+".meta")
		os.Remove(imgPath)
		os.Remove(metaPath)
		totalSize -= e.size
		c.logger.Debug("evicted cache entry", "key", e.diskKey, "freed", e.size)
	}
}

// touchMeta updates the access time in a sidecar .meta file.
func (c *ThumbnailCache) touchMeta(metaPath string, size int64) {
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		return
	}
	var m diskMeta
	if err := json.Unmarshal(raw, &m); err != nil {
		return
	}
	m.AccessTime = time.Now().Unix()
	m.Size = size
	updated, _ := json.Marshal(m)
	os.WriteFile(metaPath, updated, 0o644)
}

// MemoryLen returns the number of entries in the memory cache (for metrics).
func (c *ThumbnailCache) MemoryLen() int {
	return c.mem.Len()
}
