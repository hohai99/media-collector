package core

import (
	"fmt"
	"path/filepath"

	"media-collector/domain"
	"media-collector/filesystem"
	"media-collector/repository"
)

// MediaService handles media scanning, importing, and organisation.
type MediaService struct {
	repo    *repository.MediaRepository
	scanner *filesystem.Scanner
	mover   *filesystem.Mover
	collRepo *repository.CollectionRepository
}

// NewMediaService creates a new MediaService.
func NewMediaService(
	repo *repository.MediaRepository,
	scanner *filesystem.Scanner,
	mover *filesystem.Mover,
	collRepo *repository.CollectionRepository,
) *MediaService {
	return &MediaService{repo: repo, scanner: scanner, mover: mover, collRepo: collRepo}
}

// ScanFolder scans a directory for media files and upserts them into the DB.
// It also assigns collection_id to each media based on its directory path.
// Returns the list of discovered media.
func (s *MediaService) ScanFolder(path string) ([]domain.Media, error) {
	media, err := s.scanner.ScanFolder(path)
	if err != nil {
		return nil, fmt.Errorf("scan folder: %w", err)
	}

	for i, m := range media {
		// Match directory to a collection
		dir := filepath.ToSlash(filepath.Dir(filepath.FromSlash(m.Path)))
		coll, err := s.collRepo.GetByPath(dir)
		if err == nil && coll != nil {
			media[i].CollectionID = coll.ID
			m.CollectionID = coll.ID
		}

		if err := s.repo.Insert(m); err != nil {
			return nil, fmt.Errorf("insert media %q: %w", m.Name, err)
		}
	}

	return media, nil
}

// MoveMedia moves the given media files into the target collection's directory
// and updates the database records accordingly.
func (s *MediaService) MoveMedia(mediaIDs []string, targetCollectionPath string, targetCollectionID string) error {
	mediaItems, err := s.repo.GetByIDs(mediaIDs)
	if err != nil {
		return fmt.Errorf("get media by ids: %w", err)
	}

	for _, m := range mediaItems {
		srcPath := filepath.FromSlash(m.Path)
		dstPath := filepath.Join(filepath.FromSlash(targetCollectionPath), m.Name)

		newPath, err := s.mover.MoveFile(srcPath, dstPath)
		if err != nil {
			return fmt.Errorf("move file %q: %w", m.Name, err)
		}

		newName := filepath.Base(newPath)
		if err := s.repo.UpdatePath(m.ID, newPath, newName); err != nil {
			return fmt.Errorf("update path for %q: %w", m.Name, err)
		}
		if err := s.repo.UpdateCollection(m.ID, targetCollectionID); err != nil {
			return fmt.Errorf("update collection for %q: %w", m.Name, err)
		}
	}

	return nil
}

// GetMediaByCollection returns all media in a given collection.
func (s *MediaService) GetMediaByCollection(collectionID string) ([]domain.Media, error) {
	return s.repo.GetByCollection(collectionID)
}

// GetMediaByCollectionRecursive returns media from a collection and all its descendants.
func (s *MediaService) GetMediaByCollectionRecursive(collectionID string) ([]domain.Media, error) {
	// Get direct media
	media, err := s.repo.GetByCollection(collectionID)
	if err != nil {
		return nil, err
	}

	// Get all descendant collections
	descendants, err := s.collRepo.GetAllDescendants(collectionID)
	if err != nil {
		return nil, err
	}

	for _, desc := range descendants {
		descMedia, err := s.repo.GetByCollection(desc.ID)
		if err != nil {
			return nil, err
		}
		media = append(media, descMedia...)
	}

	return media, nil
}

// GetAllMedia returns every media record in the database.
func (s *MediaService) GetAllMedia() ([]domain.Media, error) {
	return s.repo.GetAll()
}

// GetUncategorized returns media not assigned to any collection.
func (s *MediaService) GetUncategorized() ([]domain.Media, error) {
	return s.repo.GetUncategorized()
}

