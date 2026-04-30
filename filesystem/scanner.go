package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"media-collector/domain"
)

// Scanner walks directories to discover media files.
type Scanner struct{}

// NewScanner creates a new Scanner.
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanFolder recursively walks the given path and returns all discovered media
// files. Hidden files/directories (prefixed with '.') are skipped.
func (s *Scanner) ScanFolder(root string) ([]domain.Media, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat root %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", root)
	}

	var media []domain.Media

	err = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}

		name := fi.Name()

		// Skip hidden files/directories
		if strings.HasPrefix(name, ".") {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories themselves
		if fi.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(name))
		mediaType := domain.MediaTypeFromExtension(ext)
		if mediaType == "" {
			return nil // not a recognised media file
		}

		media = append(media, domain.Media{
			ID:        uuid.New().String(),
			Name:      name,
			Path:      filepath.ToSlash(path), // normalise to forward slashes
			Type:      mediaType,
			Size:      fi.Size(),
			CreatedAt: fi.ModTime(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk %q: %w", root, err)
	}

	return media, nil
}

// ScanFolderDirect returns only the immediate (non-recursive) media files
// inside a single directory.
func (s *Scanner) ScanFolderDirect(dir string) ([]domain.Media, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", dir, err)
	}

	var media []domain.Media
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		mediaType := domain.MediaTypeFromExtension(ext)
		if mediaType == "" {
			continue
		}
		fi, err := entry.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(dir, entry.Name())
		media = append(media, domain.Media{
			ID:        uuid.New().String(),
			Name:      entry.Name(),
			Path:      filepath.ToSlash(fullPath),
			Type:      mediaType,
			Size:      fi.Size(),
			CreatedAt: fi.ModTime(),
		})
	}

	return media, nil
}

// ListSubfolders returns a list of direct subdirectory names in the given path.
func (s *Scanner) ListSubfolders(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", dir, err)
	}

	var folders []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			folders = append(folders, entry.Name())
		}
	}
	return folders, nil
}

// GetFileModTime returns the modification time of a file.
func GetFileModTime(path string) (time.Time, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fi.ModTime(), nil
}
