package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// Mover handles moving and organising files on disk.
type Mover struct{}

// NewMover creates a new Mover.
func NewMover() *Mover {
	return &Mover{}
}

// MoveFile moves src to dst. If a file already exists at dst, a numeric suffix
// is appended to avoid collision (e.g. photo_1.jpg, photo_2.jpg).
func (m *Mover) MoveFile(src, dst string) (string, error) {
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return "", err
	}

	finalPath := dst
	if _, err := os.Stat(dst); err == nil {
		// File exists — resolve conflict
		finalPath = resolveConflict(dst)
	}

	if err := os.Rename(src, finalPath); err != nil {
		return "", fmt.Errorf("move %q → %q: %w", src, finalPath, err)
	}
	return filepath.ToSlash(finalPath), nil
}

// EnsureDir creates a directory (and parents) if it does not exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// DeleteEmptyDir removes a directory only if it is empty.
func DeleteEmptyDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return nil // not empty, skip
	}
	return os.Remove(path)
}

// resolveConflict appends _1, _2, etc. to the filename until it finds a free name.
func resolveConflict(path string) string {
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
