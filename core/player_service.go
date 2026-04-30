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

// CreateConfig creates a new player config. The transition time is
// automatically calculated from the total play time and the number of
// resolved media items.
func (s *PlayerService) CreateConfig(input domain.CreateConfigInput) (*domain.PlayerConfig, error) {
	// Resolve the total media count for transition calculation
	mediaCount := len(input.MediaIDs)
	for _, colID := range input.Collections {
		// Count media in this collection + all descendants
		items, err := s.getMediaRecursive(colID)
		if err != nil {
			return nil, fmt.Errorf("count media in collection %q: %w", colID, err)
		}
		mediaCount += len(items)
	}

	if mediaCount == 0 {
		return nil, fmt.Errorf("no media selected")
	}

	transition := CalculateTransition(input.TotalPlayTime, mediaCount)

	cfg := domain.PlayerConfig{
		ID:             uuid.New().String(),
		Name:           input.Name,
		Collections:    input.Collections,
		MediaIDs:       input.MediaIDs,
		TotalPlayTime:  input.TotalPlayTime,
		TransitionTime: transition,
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

