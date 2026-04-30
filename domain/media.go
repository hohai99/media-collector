package domain

import "time"

// Media type constants
const (
	MediaTypeImage = "image"
	MediaTypeVideo = "video"
	MediaTypeAudio = "audio"
)

// Media represents a single media file on disk.
type Media struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"createdAt"`
	CollectionID string    `json:"collectionId,omitempty"`
}

// SupportedImageExtensions lists recognised image file extensions.
var SupportedImageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".bmp": true, ".webp": true, ".svg": true, ".ico": true,
}

// SupportedVideoExtensions lists recognised video file extensions.
var SupportedVideoExtensions = map[string]bool{
	".mp4": true, ".avi": true, ".mkv": true, ".mov": true,
	".wmv": true, ".flv": true, ".webm": true,
}

// SupportedAudioExtensions lists recognised audio file extensions.
var SupportedAudioExtensions = map[string]bool{
	".mp3": true, ".wav": true, ".flac": true, ".aac": true,
	".ogg": true, ".wma": true, ".m4a": true,
}

// MediaTypeFromExtension determines the media type from a file extension.
// Returns empty string if the extension is not recognised.
func MediaTypeFromExtension(ext string) string {
	if SupportedImageExtensions[ext] {
		return MediaTypeImage
	}
	if SupportedVideoExtensions[ext] {
		return MediaTypeVideo
	}
	if SupportedAudioExtensions[ext] {
		return MediaTypeAudio
	}
	return ""
}
