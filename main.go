package main

import (
	"embed"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// FileLoader serves local filesystem files to the WebView2 frontend.
// Requests to /localfile/<absolute-path> will be served from disk.
type FileLoader struct{}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestedPath := req.URL.Path

	// Only handle /localfile/ prefix
	if !strings.HasPrefix(requestedPath, "/localfile/") {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// Strip the /localfile/ prefix to get the raw URL path
	rawPath := strings.TrimPrefix(requestedPath, "/localfile/")

	// URL-decode the path (handles %20 for spaces, %3A for colons, etc.)
	filePath, err := url.PathUnescape(rawPath)
	if err != nil {
		filePath = rawPath
	}

	// Convert to OS-specific path separators
	filePath = filepath.FromSlash(filePath)

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// Use http.ServeFile to stream the file instead of loading it all into
	// memory. This prevents OOM crashes on large video files (e.g. 1.7 GB)
	// and enables HTTP Range requests for video seeking.
	http.ServeFile(res, req, filePath)
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
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

	if err != nil {
		println("Error:", err.Error())
	}
}

