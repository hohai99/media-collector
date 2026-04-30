package domain

// PlayerConfig defines a playback preset with timing and media selection.
type PlayerConfig struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Collections    []string `json:"collections,omitempty"`
	MediaIDs       []string `json:"mediaIds,omitempty"`
	TotalPlayTime  int      `json:"totalPlayTime"`  // seconds
	TransitionTime int      `json:"transitionTime"` // seconds
}

// CreateConfigInput is the input payload for creating a new player config.
type CreateConfigInput struct {
	Name          string   `json:"name"`
	Collections   []string `json:"collections,omitempty"`
	MediaIDs      []string `json:"mediaIds,omitempty"`
	TotalPlayTime int      `json:"totalPlayTime"`
}
