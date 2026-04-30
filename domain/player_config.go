package domain

// PlayerConfig defines a playback preset with timing and media selection.
type PlayerConfig struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Collections    []string `json:"collections,omitempty"`
	MediaIDs       []string `json:"mediaIds,omitempty"`
	TimePerPicture int      `json:"timePerPicture"` // seconds per image slide
	TotalPlayTime  int      `json:"totalPlayTime"`  // seconds (auto-calculated)
	TransitionTime int      `json:"transitionTime"` // seconds (kept for compat, = TimePerPicture)
	SoundSource    string   `json:"soundSource,omitempty"` // optional background audio path
}

// CreateConfigInput is the input payload for creating a new player config.
type CreateConfigInput struct {
	Name           string   `json:"name"`
	Collections    []string `json:"collections,omitempty"`
	MediaIDs       []string `json:"mediaIds,omitempty"`
	TimePerPicture int      `json:"timePerPicture"` // seconds per image slide
	SoundSource    string   `json:"soundSource,omitempty"`
}
