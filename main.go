package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"media-collector/server"
	"media-collector/thumbs"
	"media-collector/transcode"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Structured logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Thumbnail generator (configurable cache sizes)
	thumbCfg := thumbs.DefaultConfig()
	thumbGen, err := thumbs.NewGenerator(thumbCfg, logger)
	if err != nil {
		logger.Error("failed to init thumbnail generator", "error", err)
		os.Exit(1)
	}

	// Video transcoder (optional — requires ffmpeg binary)
	var tc *transcode.Transcoder
	tcCfg := transcode.DefaultConfig()
	if err := tcCfg.ValidateFFmpeg(); err != nil {
		logger.Warn("ffmpeg not available, video transcoding disabled", "error", err)
	} else {
		tc, err = transcode.NewTranscoder(tcCfg, logger)
		if err != nil {
			logger.Warn("failed to init transcoder", "error", err)
		} else {
			logger.Info("video transcoding enabled", "ffmpeg", tcCfg.FFmpegPath)
		}
	}

	// Standalone file server
	fileServer, err := server.New(thumbGen, tc, logger)
	if err != nil {
		logger.Error("failed to init file server", "error", err)
		os.Exit(1)
	}
	fileServer.Start()
	logger.Info("file server ready", "port", fileServer.Port())

	app := NewApp(fileServer, logger)

	wailsErr := wails.Run(&options.App{
		Title:     "Media Collector",
		Width:     1280,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 15, G: 17, B: 26, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if wailsErr != nil {
		logger.Error("wails error", "error", wailsErr)
	}
}

// FileLoader is kept as a thin fallback for the Wails asset pipeline.
type FileLoader struct{}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}
