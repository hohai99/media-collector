package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFolder_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	scanner := NewScanner()

	media, err := scanner.ScanFolder(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(media) != 0 {
		t.Errorf("expected 0 media, got %d", len(media))
	}
}

func TestScanFolder_DetectsMediaTypes(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	files := []string{"photo.jpg", "video.mp4", "song.mp3", "readme.txt", "code.go"}
	for _, f := range files {
		os.WriteFile(filepath.Join(dir, f), []byte("test"), 0o644)
	}

	scanner := NewScanner()
	media, err := scanner.ScanFolder(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(media) != 3 {
		t.Errorf("expected 3 media files, got %d", len(media))
	}

	typeMap := map[string]bool{}
	for _, m := range media {
		typeMap[m.Type] = true
	}
	if !typeMap["image"] {
		t.Error("expected to find image type")
	}
	if !typeMap["video"] {
		t.Error("expected to find video type")
	}
	if !typeMap["audio"] {
		t.Error("expected to find audio type")
	}
}

func TestScanFolder_SkipsHiddenFiles(t *testing.T) {
	dir := t.TempDir()

	// Create a hidden directory
	hiddenDir := filepath.Join(dir, ".hidden")
	os.MkdirAll(hiddenDir, 0o755)
	os.WriteFile(filepath.Join(hiddenDir, "secret.jpg"), []byte("test"), 0o644)

	// Create a hidden file
	os.WriteFile(filepath.Join(dir, ".hidden.png"), []byte("test"), 0o644)

	// Create a normal file
	os.WriteFile(filepath.Join(dir, "visible.png"), []byte("test"), 0o644)

	scanner := NewScanner()
	media, err := scanner.ScanFolder(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(media) != 1 {
		t.Errorf("expected 1 media file (skipping hidden), got %d", len(media))
	}
	if len(media) > 0 && media[0].Name != "visible.png" {
		t.Errorf("expected visible.png, got %s", media[0].Name)
	}
}

func TestScanFolder_NestedFolders(t *testing.T) {
	dir := t.TempDir()

	// Create nested structure
	sub := filepath.Join(dir, "sub", "deep")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(dir, "root.jpg"), []byte("test"), 0o644)
	os.WriteFile(filepath.Join(dir, "sub", "mid.png"), []byte("test"), 0o644)
	os.WriteFile(filepath.Join(sub, "deep.mp4"), []byte("test"), 0o644)

	scanner := NewScanner()
	media, err := scanner.ScanFolder(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(media) != 3 {
		t.Errorf("expected 3 media files from nested dirs, got %d", len(media))
	}
}

func TestScanFolder_InvalidPath(t *testing.T) {
	scanner := NewScanner()
	_, err := scanner.ScanFolder("/nonexistent/path/abc123")
	if err == nil {
		t.Error("expected error for nonexistent path")
	}
}
