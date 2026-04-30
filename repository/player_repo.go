package repository

import (
	"database/sql"
	"fmt"

	"media-collector/domain"
)

// PlayerRepository provides CRUD access to player configs and their media/collection associations.
type PlayerRepository struct {
	db *sql.DB
}

// NewPlayerRepository creates a new PlayerRepository.
func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

// InsertConfig creates a new player config and its media/collection associations.
func (r *PlayerRepository) InsertConfig(cfg domain.PlayerConfig) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO player_configs (id, name, total_time, transition_time, time_per_picture, sound_source) VALUES (?, ?, ?, ?, ?, ?)`,
		cfg.ID, cfg.Name, cfg.TotalPlayTime, cfg.TransitionTime, cfg.TimePerPicture, cfg.SoundSource,
	)
	if err != nil {
		return fmt.Errorf("insert player config: %w", err)
	}

	// Insert media associations
	for i, mediaID := range cfg.MediaIDs {
		_, err = tx.Exec(
			`INSERT INTO player_media (config_id, media_id, position) VALUES (?, ?, ?)`,
			cfg.ID, mediaID, i,
		)
		if err != nil {
			return fmt.Errorf("insert player media: %w", err)
		}
	}

	// Insert collection associations
	for _, colID := range cfg.Collections {
		_, err = tx.Exec(
			`INSERT INTO player_collections (config_id, collection_id) VALUES (?, ?)`,
			cfg.ID, colID,
		)
		if err != nil {
			return fmt.Errorf("insert player collection: %w", err)
		}
	}

	return tx.Commit()
}

// GetConfig retrieves a player config by ID, including its media and collection IDs.
func (r *PlayerRepository) GetConfig(id string) (*domain.PlayerConfig, error) {
	var cfg domain.PlayerConfig
	err := r.db.QueryRow(
		`SELECT id, name, total_time, transition_time, time_per_picture, sound_source FROM player_configs WHERE id = ?`, id,
	).Scan(&cfg.ID, &cfg.Name, &cfg.TotalPlayTime, &cfg.TransitionTime, &cfg.TimePerPicture, &cfg.SoundSource)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query player config: %w", err)
	}

	// Load media IDs
	mediaRows, err := r.db.Query(
		`SELECT media_id FROM player_media WHERE config_id = ? ORDER BY position`, id)
	if err != nil {
		return nil, fmt.Errorf("query player media: %w", err)
	}
	defer mediaRows.Close()
	for mediaRows.Next() {
		var mid string
		if err := mediaRows.Scan(&mid); err != nil {
			return nil, fmt.Errorf("scan player media: %w", err)
		}
		cfg.MediaIDs = append(cfg.MediaIDs, mid)
	}

	// Load collection IDs
	colRows, err := r.db.Query(
		`SELECT collection_id FROM player_collections WHERE config_id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("query player collections: %w", err)
	}
	defer colRows.Close()
	for colRows.Next() {
		var cid string
		if err := colRows.Scan(&cid); err != nil {
			return nil, fmt.Errorf("scan player collection: %w", err)
		}
		cfg.Collections = append(cfg.Collections, cid)
	}

	return &cfg, nil
}

// GetAllConfigs returns every player config (without media/collection details for listing).
func (r *PlayerRepository) GetAllConfigs() ([]domain.PlayerConfig, error) {
	rows, err := r.db.Query(
		`SELECT id, name, total_time, transition_time, time_per_picture, sound_source FROM player_configs ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query all player configs: %w", err)
	}
	defer rows.Close()

	var results []domain.PlayerConfig
	for rows.Next() {
		var cfg domain.PlayerConfig
		if err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.TotalPlayTime, &cfg.TransitionTime, &cfg.TimePerPicture, &cfg.SoundSource); err != nil {
			return nil, fmt.Errorf("scan player config: %w", err)
		}
		results = append(results, cfg)
	}
	return results, rows.Err()
}

// DeleteConfig removes a player config and cascades to its associations.
func (r *PlayerRepository) DeleteConfig(id string) error {
	_, err := r.db.Exec(`DELETE FROM player_configs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete player config: %w", err)
	}
	return nil
}

// GetMediaIDsForConfig returns all media IDs associated with a config.
func (r *PlayerRepository) GetMediaIDsForConfig(configID string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT media_id FROM player_media WHERE config_id = ? ORDER BY position`, configID)
	if err != nil {
		return nil, fmt.Errorf("query config media ids: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan media id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// GetCollectionIDsForConfig returns all collection IDs associated with a config.
func (r *PlayerRepository) GetCollectionIDsForConfig(configID string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT collection_id FROM player_collections WHERE config_id = ?`, configID)
	if err != nil {
		return nil, fmt.Errorf("query config collection ids: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan collection id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
