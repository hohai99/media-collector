package repository

import (
	"database/sql"
	"fmt"

	"media-collector/domain"
)

// MediaRepository provides CRUD access to the media table.
type MediaRepository struct {
	db *sql.DB
}

// NewMediaRepository creates a new MediaRepository.
func NewMediaRepository(db *sql.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

// Insert adds a new media record to the database.
func (r *MediaRepository) Insert(m domain.Media) error {
	_, err := r.db.Exec(
		`INSERT INTO media (id, name, path, type, size, created_at, collection_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET name=excluded.name, size=excluded.size, type=excluded.type`,
		m.ID, m.Name, m.Path, m.Type, m.Size, m.CreatedAt, nilIfEmpty(m.CollectionID),
	)
	if err != nil {
		return fmt.Errorf("insert media: %w", err)
	}
	return nil
}

// GetAll returns every media record.
func (r *MediaRepository) GetAll() ([]domain.Media, error) {
	rows, err := r.db.Query(
		`SELECT id, name, path, type, size, created_at, COALESCE(collection_id, '') FROM media ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query all media: %w", err)
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// GetByCollection returns media belonging to a specific collection.
func (r *MediaRepository) GetByCollection(collectionID string) ([]domain.Media, error) {
	rows, err := r.db.Query(
		`SELECT id, name, path, type, size, created_at, COALESCE(collection_id, '')
		 FROM media WHERE collection_id = ? ORDER BY name`, collectionID)
	if err != nil {
		return nil, fmt.Errorf("query media by collection: %w", err)
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// GetByIDs returns media matching the given IDs.
func (r *MediaRepository) GetByIDs(ids []string) ([]domain.Media, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := `SELECT id, name, path, type, size, created_at, COALESCE(collection_id, '') FROM media WHERE id IN (`
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = id
	}
	query += `) ORDER BY name`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query media by ids: %w", err)
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

// UpdateCollection moves a media item to a different collection.
func (r *MediaRepository) UpdateCollection(mediaID, collectionID string) error {
	_, err := r.db.Exec(`UPDATE media SET collection_id = ? WHERE id = ?`,
		nilIfEmpty(collectionID), mediaID)
	if err != nil {
		return fmt.Errorf("update media collection: %w", err)
	}
	return nil
}

// UpdatePath updates the file path of a media record (e.g. after a move).
func (r *MediaRepository) UpdatePath(mediaID, newPath, newName string) error {
	_, err := r.db.Exec(`UPDATE media SET path = ?, name = ? WHERE id = ?`,
		newPath, newName, mediaID)
	if err != nil {
		return fmt.Errorf("update media path: %w", err)
	}
	return nil
}

// Delete removes a media record by ID.
func (r *MediaRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM media WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete media: %w", err)
	}
	return nil
}

// GetUncategorized returns media that don't belong to any collection.
func (r *MediaRepository) GetUncategorized() ([]domain.Media, error) {
	rows, err := r.db.Query(
		`SELECT id, name, path, type, size, created_at, '' FROM media WHERE collection_id IS NULL ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query uncategorized media: %w", err)
	}
	defer rows.Close()
	return scanMediaRows(rows)
}

func scanMediaRows(rows *sql.Rows) ([]domain.Media, error) {
	var results []domain.Media
	for rows.Next() {
		var m domain.Media
		if err := rows.Scan(&m.ID, &m.Name, &m.Path, &m.Type, &m.Size, &m.CreatedAt, &m.CollectionID); err != nil {
			return nil, fmt.Errorf("scan media row: %w", err)
		}
		results = append(results, m)
	}
	return results, rows.Err()
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
