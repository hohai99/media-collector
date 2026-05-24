package transcode

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds all configurable settings for the video transcoding subsystem.
type Config struct {
	// FFmpegPath is the absolute path to the ffmpeg binary.
	// Auto-resolved from {exeDir}/ffmpeg.exe if empty.
	FFmpegPath string

	// CacheDir is the directory for HLS segment output.
	// Defaults to {exeDir}/cache/hls if empty.
	CacheDir string

	// SegmentSecs is the HLS segment duration in seconds.
	SegmentSecs int

	// MaxConcurrent is the maximum number of concurrent transcode jobs.
	MaxConcurrent int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		FFmpegPath:    resolveFFmpegPath(),
		CacheDir:      defaultHLSCacheDir(),
		SegmentSecs:   4,
		MaxConcurrent: 2,
	}
}

// resolveFFmpegPath looks for ffmpeg next to the running executable.
func resolveFFmpegPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ffmpegBinaryName()
	}
	dir := filepath.Dir(exe)
	candidate := filepath.Join(dir, ffmpegBinaryName())
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	// Fall back to PATH lookup
	return ffmpegBinaryName()
}

// ffmpegBinaryName returns the platform-specific ffmpeg binary name.
func ffmpegBinaryName() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

// defaultHLSCacheDir resolves {exeDir}/cache/hls.
func defaultHLSCacheDir() string {
	exe, err := os.Executable()
	if err != nil {
		return filepath.Join(".", "cache", "hls")
	}
	return filepath.Join(filepath.Dir(exe), "cache", "hls")
}

// ValidateFFmpeg checks whether ffmpeg is available at the configured path.
func (c Config) ValidateFFmpeg() error {
	if _, err := os.Stat(c.FFmpegPath); err != nil {
		return fmt.Errorf("ffmpeg not found at %q: %w", c.FFmpegPath, err)
	}
	return nil
}
