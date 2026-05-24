package transcode

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Job represents an active or completed transcode job.
type Job struct {
	ID         string // sha256 hash of source path
	SourcePath string
	OutputDir  string // directory containing m3u8 + ts segments

	mu   sync.Mutex
	cmd  *exec.Cmd
	done chan struct{}
	err  error
}

// Wait blocks until the transcode job has completed (or failed).
func (j *Job) Wait() error {
	<-j.done
	return j.err
}

// PlaylistPath returns the path to the HLS playlist file.
func (j *Job) PlaylistPath() string {
	return filepath.Join(j.OutputDir, "playlist.m3u8")
}

// Transcoder manages ffmpeg-based HLS transcoding with concurrency control.
type Transcoder struct {
	cfg    Config
	jobs   map[string]*Job
	mu     sync.RWMutex
	sem    chan struct{}
	logger *slog.Logger
}

// NewTranscoder creates a Transcoder. It cleans up any leftover HLS cache
// from a previous app session.
func NewTranscoder(cfg Config, logger *slog.Logger) (*Transcoder, error) {
	// Wipe HLS cache on startup (per requirements)
	if err := os.RemoveAll(cfg.CacheDir); err != nil {
		logger.Warn("failed to clean HLS cache", "dir", cfg.CacheDir, "error", err)
	}
	if err := os.MkdirAll(cfg.CacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("create HLS cache dir: %w", err)
	}

	return &Transcoder{
		cfg:    cfg,
		jobs:   make(map[string]*Job),
		sem:    make(chan struct{}, cfg.MaxConcurrent),
		logger: logger,
	}, nil
}

// hashPath produces a short deterministic ID for a video file path.
func hashPath(path string) string {
	sum := sha256.Sum256([]byte(path))
	return fmt.Sprintf("%x", sum[:16]) // 32 hex chars
}

// EnsureReady returns a Job for the given video file. If the file has already
// been transcoded (segments exist), it returns immediately. Otherwise it kicks
// off an ffmpeg process and waits until at least the playlist file is created.
func (t *Transcoder) EnsureReady(ctx context.Context, videoPath string) (*Job, error) {
	id := hashPath(videoPath)
	outputDir := filepath.Join(t.cfg.CacheDir, id)
	playlistPath := filepath.Join(outputDir, "playlist.m3u8")

	// Fast path: existing job
	t.mu.RLock()
	if job, ok := t.jobs[id]; ok {
		t.mu.RUnlock()
		select {
		case <-job.done:
			if job.err != nil {
				break // fall through to re-transcode
			}
			return job, nil
		default:
			// Job in progress — wait for playlist file
			return job, t.waitForPlaylist(ctx, playlistPath)
		}
	} else {
		t.mu.RUnlock()
	}

	// Check if segments are already cached on disk
	if _, err := os.Stat(playlistPath); err == nil {
		job := &Job{
			ID:         id,
			SourcePath: videoPath,
			OutputDir:  outputDir,
			done:       make(chan struct{}),
		}
		close(job.done) // already complete
		t.mu.Lock()
		t.jobs[id] = job
		t.mu.Unlock()
		return job, nil
	}

	// Need to transcode
	return t.startJob(ctx, id, videoPath, outputDir)
}

func (t *Transcoder) startJob(ctx context.Context, id, videoPath, outputDir string) (*Job, error) {
	// Acquire semaphore
	select {
	case t.sem <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		<-t.sem
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	playlistPath := filepath.Join(outputDir, "playlist.m3u8")
	segPattern := filepath.Join(outputDir, "seg_%03d.ts")

	cmd := exec.CommandContext(ctx, t.cfg.FFmpegPath,
		"-i", videoPath,
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "hls",
		"-hls_time", strconv.Itoa(t.cfg.SegmentSecs),
		"-hls_list_size", "0",
		"-hls_segment_filename", segPattern,
		playlistPath,
	)

	job := &Job{
		ID:         id,
		SourcePath: videoPath,
		OutputDir:  outputDir,
		cmd:        cmd,
		done:       make(chan struct{}),
	}

	t.mu.Lock()
	t.jobs[id] = job
	t.mu.Unlock()

	t.logger.Info("starting transcode", "id", id, "source", videoPath)

	if err := cmd.Start(); err != nil {
		<-t.sem
		job.err = fmt.Errorf("start ffmpeg: %w", err)
		close(job.done)
		return nil, job.err
	}

	// Wait for completion in background
	go func() {
		defer func() {
			<-t.sem
			close(job.done)
		}()
		if err := cmd.Wait(); err != nil {
			job.mu.Lock()
			job.err = fmt.Errorf("ffmpeg: %w", err)
			job.mu.Unlock()
			t.logger.Warn("transcode failed", "id", id, "error", err)
		} else {
			t.logger.Info("transcode complete", "id", id)
		}
	}()

	// Wait until the playlist file appears (or context cancelled)
	if err := t.waitForPlaylist(ctx, playlistPath); err != nil {
		return nil, err
	}

	return job, nil
}

// waitForPlaylist polls for the playlist file to appear on disk.
func (t *Transcoder) waitForPlaylist(ctx context.Context, path string) error {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// continue polling
		}
	}
}

// Shutdown kills all running ffmpeg processes and cleans up.
func (t *Transcoder) Shutdown() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, job := range t.jobs {
		select {
		case <-job.done:
			// Already finished
		default:
			if job.cmd != nil && job.cmd.Process != nil {
				t.logger.Info("killing transcode process", "id", job.ID)
				job.cmd.Process.Kill()
			}
		}
	}
}

// Cleanup removes all HLS cached segments.
func (t *Transcoder) Cleanup() error {
	return os.RemoveAll(t.cfg.CacheDir)
}
