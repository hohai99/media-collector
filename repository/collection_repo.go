package repository

import (
	"database/sql"
	"fmt"

	"media-collector/domain"
)

// CollectionRepository provides CRUD access to the collections table.
type CollectionRepository struct {
	db *sql.DB
}

// NewCollectionRepository creates a new CollectionRepository.
func NewCollectionRepository(db *sql.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Insert adds a new collection record.
func (r *CollectionRepository) Insert(c domain.Collection) error {
	_, err := r.db.Exec(
		`INSERT INTO collections (id, name, path, parent_id) VALUES (?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET name=excluded.name, parent_id=excluded.parent_id`,
		c.ID, c.Name, c.Path, c.ParentID,
	)
	if err != nil {
		return fmt.Errorf("insert collection: %w", err)
	}
	return nil
}

// GetAll returns every collection.
func (r *CollectionRepository) GetAll() ([]domain.Collection, error) {
	rows, err := r.db.Query(
		`SELECT id, name, path, parent_id FROM collections ORDER BY path`)
	if err != nil {
		return nil, fmt.Errorf("query all collections: %w", err)
	}
	defer rows.Close()
	return scanCollectionRows(rows)
}

// GetByID returns a single collection by ID.
func (r *CollectionRepository) GetByID(id string) (*domain.Collection, error) {
	row := r.db.QueryRow(
		`SELECT id, name, path, parent_id FROM collections WHERE id = ?`, id)
	var c domain.Collection
	if err := row.Scan(&c.ID, &c.Name, &c.Path, &c.ParentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query collection by id: %w", err)
	}
	return &c, nil
}

// GetByPath returns a single collection by path.
func (r *CollectionRepository) GetByPath(path string) (*domain.Collection, error) {
	row := r.db.QueryRow(
		`SELECT id, name, path, parent_id FROM collections WHERE path = ?`, path)
	var c domain.Collection
	if err := row.Scan(&c.ID, &c.Name, &c.Path, &c.ParentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query collection by path: %w", err)
	}
	return &c, nil
}

// GetChildren returns direct children of a parent collection.
func (r *CollectionRepository) GetChildren(parentID string) ([]domain.Collection, error) {
	rows, err := r.db.Query(
		`SELECT id, name, path, parent_id FROM collections WHERE parent_id = ? ORDER BY name`,
		parentID)
	if err != nil {
		return nil, fmt.Errorf("query children: %w", err)
	}
	defer rows.Close()
	return scanCollectionRows(rows)
}

// Delete removes a collection by ID.
func (r *CollectionRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM collections WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

// GetAllDescendants returns all descendant collections of a parent (recursive BFS).
func (r *CollectionRepository) GetAllDescendants(parentID string) ([]domain.Collection, error) {
	var all []domain.Collection
	queue := []string{parentID}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		children, err := r.GetChildren(current)
		if err != nil {
			return nil, err
		}
		for _, ch := range children {
			all = append(all, ch)
			queue = append(queue, ch.ID)
		}
	}
	return all, nil
}

// GetTopLevel returns collections that have no parent (root-level collections).
func (r *CollectionRepository) GetTopLevel() ([]domain.Collection, error) {
	rows, err := r.db.Query(
		`SELECT id, name, path, parent_id FROM collections WHERE parent_id IS NULL ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query top-level collections: %w", err)
	}
	defer rows.Close()
	return scanCollectionRows(rows)
}

func scanCollectionRows(rows *sql.Rows) ([]domain.Collection, error) {
	var results []domain.Collection
	for rows.Next() {
		var c domain.Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.Path, &c.ParentID); err != nil {
			return nil, fmt.Errorf("scan collection row: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}
