package domain

import (
	"testing"
)

func TestMediaTypeFromExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".jpg", MediaTypeImage},
		{".jpeg", MediaTypeImage},
		{".png", MediaTypeImage},
		{".gif", MediaTypeImage},
		{".webp", MediaTypeImage},
		{".mp4", MediaTypeVideo},
		{".avi", MediaTypeVideo},
		{".mkv", MediaTypeVideo},
		{".webm", MediaTypeVideo},
		{".mp3", MediaTypeAudio},
		{".wav", MediaTypeAudio},
		{".flac", MediaTypeAudio},
		{".txt", ""},
		{".go", ""},
		{".exe", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := MediaTypeFromExtension(tt.ext)
			if result != tt.expected {
				t.Errorf("MediaTypeFromExtension(%q) = %q; want %q",
					tt.ext, result, tt.expected)
			}
		})
	}
}
