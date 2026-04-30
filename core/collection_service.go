package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"media-collector/domain"
	"media-collector/filesystem"
	"media-collector/repository"
)

// CollectionService handles creating, syncing, and querying collections.
type CollectionService struct {
	repo *repository.CollectionRepository
}

// NewCollectionService creates a new CollectionService.
func NewCollectionService(repo *repository.CollectionRepository) *CollectionService {
	return &CollectionService{repo: repo}
}

// CreateCollection creates a new collection folder on disk and in the DB.
func (s *CollectionService) CreateCollection(name string, parentID *string, masterFolder string) (*domain.Collection, error) {
	// Determine parent path
	parentPath := masterFolder
	if parentID != nil && *parentID != "" {
		parent, err := s.repo.GetByID(*parentID)
		if err != nil {
			return nil, fmt.Errorf("get parent: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent collection %q not found", *parentID)
		}
		parentPath = filepath.FromSlash(parent.Path)
	}

	// Create directory
	collPath := filepath.Join(parentPath, name)
	if err := filesystem.EnsureDir(collPath); err != nil {
		return nil, fmt.Errorf("create collection dir: %w", err)
	}

	coll := domain.Collection{
		ID:       uuid.New().String(),
		Name:     name,
		Path:     filepath.ToSlash(collPath),
		ParentID: parentID,
	}

	if err := s.repo.Insert(coll); err != nil {
		return nil, fmt.Errorf("insert collection: %w", err)
	}

	return &coll, nil
}

// GetCollectionTree returns all collections in a flat list (ordered by path).
// The frontend is responsible for building the tree from parentID references.
func (s *CollectionService) GetCollectionTree() ([]domain.Collection, error) {
	return s.repo.GetAll()
}

// DeleteCollection removes a collection from the DB. Does NOT delete the folder.
func (s *CollectionService) DeleteCollection(id string) error {
	return s.repo.Delete(id)
}

// SyncCollections walks the master folder and ensures every subfolder is
// represented in the DB. This reconciles the filesystem → DB direction.
func (s *CollectionService) SyncCollections(masterFolder string) error {
	return filepath.Walk(masterFolder, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable
		}
		if !fi.IsDir() {
			return nil
		}
		if strings.HasPrefix(fi.Name(), ".") {
			return filepath.SkipDir
		}
		// Skip the master folder itself
		if path == masterFolder {
			return nil
		}

		normalPath := filepath.ToSlash(path)

		// Check if already in DB
		existing, err := s.repo.GetByPath(normalPath)
		if err != nil {
			return fmt.Errorf("check path %q: %w", normalPath, err)
		}
		if existing != nil {
			return nil // already tracked
		}

		// Find parent
		parentDir := filepath.ToSlash(filepath.Dir(path))
		var parentID *string
		if parentDir != filepath.ToSlash(masterFolder) {
			parent, err := s.repo.GetByPath(parentDir)
			if err != nil {
				return fmt.Errorf("find parent for %q: %w", path, err)
			}
			if parent != nil {
				parentID = &parent.ID
			}
		}

		coll := domain.Collection{
			ID:       uuid.New().String(),
			Name:     fi.Name(),
			Path:     normalPath,
			ParentID: parentID,
		}
		return s.repo.Insert(coll)
	})
}

// GetThumbnail finds the first image file in the collection's directory.
// Returns empty string if no images are found.
func (s *CollectionService) GetThumbnail(collectionPath string) string {
	dir := filepath.FromSlash(collectionPath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if domain.SupportedImageExtensions[ext] {
			fullPath := filepath.Join(dir, entry.Name())
			return filepath.ToSlash(fullPath)
		}
	}
	return ""
}

// GetDescendantIDs returns the IDs of all descendant collections (recursive).
func (s *CollectionService) GetDescendantIDs(collectionID string) ([]string, error) {
	descendants, err := s.repo.GetAllDescendants(collectionID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(descendants))
	for i, d := range descendants {
		ids[i] = d.ID
	}
	return ids, nil
}

