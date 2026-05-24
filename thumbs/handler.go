package thumbs

import (
	"net/http"
	"strconv"
)

// ServeHTTP handles GET /thumb?path=...&w=...&h=...&fit=...
func (g *Generator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	q := r.URL.Query()

	// Required: path
	filePath := q.Get("path")
	if filePath == "" {
		http.Error(w, "missing 'path' parameter", http.StatusBadRequest)
		return
	}

	// Validate it's an image
	if !isImagePath(filePath) {
		http.Error(w, "unsupported image format", http.StatusBadRequest)
		return
	}

	// Optional: w, h, fit (with defaults)
	width := g.cfg.DefaultWidth
	height := g.cfg.DefaultHeight
	fit := g.cfg.DefaultFit

	if v := q.Get("w"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 2000 {
			width = parsed
		}
	}
	if v := q.Get("h"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 2000 {
			height = parsed
		}
	}
	if v := q.Get("fit"); v == "cover" || v == "contain" {
		fit = v
	}

	data, err := g.Generate(r.Context(), filePath, width, height, fit)
	if err != nil {
		if r.Context().Err() != nil {
			// Client disconnected — don't log as error.
			return
		}
		g.logger.Warn("thumbnail generation failed",
			"path", filePath,
			"error", err,
		)
		http.Error(w, "thumbnail generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "max-age=86400")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}
