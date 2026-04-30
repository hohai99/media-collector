package domain

// Collection represents a folder-based grouping of media files.
type Collection struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	ParentID *string `json:"parentId,omitempty"`
}
