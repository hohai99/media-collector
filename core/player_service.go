package core

import (
	"fmt"

	"github.com/google/uuid"

	"media-collector/domain"
	"media-collector/repository"
)

// PlayerService handles player configuration and media resolution.
type PlayerService struct {
	playerRepo *repository.PlayerRepository
	mediaRepo  *repository.MediaRepository
	collRepo   *repository.CollectionRepository
}

// NewPlayerService creates a new PlayerService.
func NewPlayerService(
	playerRepo *repository.PlayerRepository,
	mediaRepo *repository.MediaRepository,
	collRepo *repository.CollectionRepository,
) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		mediaRepo:  mediaRepo,
		collRepo:   collRepo,
	}
}

// CreateConfig creates a new player config. The total play time is
// auto-calculated: totalPlayTime = numberOfImages × timePerPicture.
// Videos play at their natural duration and are not counted.
func (s *PlayerService) CreateConfig(input domain.CreateConfigInput) (*domain.PlayerConfig, error) {
	// Resolve all media to count images
	var allMedia []domain.Media

	// Direct media
	if len(input.MediaIDs) > 0 {
		directMedia, err := s.mediaRepo.GetByIDs(input.MediaIDs)
		if err != nil {
			return nil, fmt.Errorf("get direct media: %w", err)
		}
		allMedia = append(allMedia, directMedia...)
	}

	// Collection media (recursive)
	for _, colID := range input.Collections {
		colMedia, err := s.getMediaRecursive(colID)
		if err != nil {
			return nil, fmt.Errorf("get collection %q media: %w", colID, err)
		}
		allMedia = append(allMedia, colMedia...)
	}

	if len(allMedia) == 0 {
		return nil, fmt.Errorf("no media selected")
	}

	// Count images for time calculation
	imageCount := 0
	for _, m := range allMedia {
		if m.Type == domain.MediaTypeImage {
			imageCount++
		}
	}

	timePerPicture := input.TimePerPicture
	if timePerPicture <= 0 {
		timePerPicture = 5 // default 5 seconds
	}

	// Auto-calculate total play time for images only
	// (videos play at their natural duration, not factored in)
	totalPlayTime := imageCount * timePerPicture

	cfg := domain.PlayerConfig{
		ID:             uuid.New().String(),
		Name:           input.Name,
		Collections:    input.Collections,
		MediaIDs:       input.MediaIDs,
		TimePerPicture: timePerPicture,
		TotalPlayTime:  totalPlayTime,
		TransitionTime: timePerPicture, // kept for backward compat
		SoundSource:    input.SoundSource,
	}

	if err := s.playerRepo.InsertConfig(cfg); err != nil {
		return nil, fmt.Errorf("insert player config: %w", err)
	}

	return &cfg, nil
}

// GetConfig retrieves a player config by ID.
func (s *PlayerService) GetConfig(id string) (*domain.PlayerConfig, error) {
	return s.playerRepo.GetConfig(id)
}

// GetAllConfigs returns all player configs.
func (s *PlayerService) GetAllConfigs() ([]domain.PlayerConfig, error) {
	return s.playerRepo.GetAllConfigs()
}

// DeleteConfig removes a player config.
func (s *PlayerService) DeleteConfig(id string) error {
	return s.playerRepo.DeleteConfig(id)
}

// CalculateTransition computes seconds per media item.
// transition = totalTime / numberOfItems
func CalculateTransition(totalTime int, mediaCount int) int {
	if mediaCount <= 0 {
		return 0
	}
	return totalTime / mediaCount
}

// ResolveMediaList gathers all media for a player config by merging
// direct media IDs and media from associated collections (recursively).
// Duplicates are removed (direct IDs take priority).
func (s *PlayerService) ResolveMediaList(configID string) ([]domain.Media, error) {
	cfg, err := s.playerRepo.GetConfig(configID)
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("config %q not found", configID)
	}

	seen := make(map[string]bool)
	var result []domain.Media

	// Direct media first
	if len(cfg.MediaIDs) > 0 {
		directMedia, err := s.mediaRepo.GetByIDs(cfg.MediaIDs)
		if err != nil {
			return nil, fmt.Errorf("get direct media: %w", err)
		}
		for _, m := range directMedia {
			if !seen[m.ID] {
				seen[m.ID] = true
				result = append(result, m)
			}
		}
	}

	// Then collection media (recursive — includes sub-collections)
	for _, colID := range cfg.Collections {
		colMedia, err := s.getMediaRecursive(colID)
		if err != nil {
			return nil, fmt.Errorf("get collection %q media: %w", colID, err)
		}
		for _, m := range colMedia {
			if !seen[m.ID] {
				seen[m.ID] = true
				result = append(result, m)
			}
		}
	}

	return result, nil
}

// getMediaRecursive returns media from a collection and all its descendants.
func (s *PlayerService) getMediaRecursive(collectionID string) ([]domain.Media, error) {
	// Direct media in this collection
	media, err := s.mediaRepo.GetByCollection(collectionID)
	if err != nil {
		return nil, err
	}

	// Descendant collections
	descendants, err := s.collRepo.GetAllDescendants(collectionID)
	if err != nil {
		return nil, err
	}
	for _, desc := range descendants {
		descMedia, err := s.mediaRepo.GetByCollection(desc.ID)
		if err != nil {
			return nil, err
		}
		media = append(media, descMedia...)
	}

	return media, nil
}

