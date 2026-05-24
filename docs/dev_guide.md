# Media Collector — Developer Guide

## Prerequisites

| Tool | Version / Requirement | Install / Source |
|------|-----------------------|------------------|
| Go | ≥ 1.21 | https://go.dev/dl/ |
| Node.js | ≥ 18 | https://nodejs.org/ |
| Wails CLI | ≥ 2.x | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| FFmpeg | Optional | https://ffmpeg.org/download.html (Required for video HLS transcoding) |

## Getting Started

```bash
# Clone the repo
git clone <repo-url>
cd media-collector

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run in development mode (hot-reload)
wails dev
```

## Setting up FFmpeg (Optional)

The application utilizes **FFmpeg** to transcode videos into HLS segments on-the-fly for smooth web playback and seeking.

### How to Install & Configure

1. **Local Binary (Portable / Recommended)**:
   - Download the static build of FFmpeg for your OS from the [official FFmpeg site](https://ffmpeg.org/download.html).
   - Place the executable in the **same directory as the application executable**:
     - **During Development (`wails dev`)**: Place `ffmpeg.exe` (Windows) or `ffmpeg` (macOS/Linux) inside the **`build/bin/`** directory.
     - **During Production (Built App)**: Place `ffmpeg.exe` next to the compiled `Media Collector.exe` in your distribution folder.
2. **System PATH**:
   - Alternatively, install FFmpeg system-wide (e.g., via `winget install Gyan.FFmpeg` on Windows or `brew install ffmpeg` on macOS) and ensure it is available in your system environment `PATH`.

### Verification

Upon startup, the backend automatically scans for `ffmpeg` next to the executable, falling back to looking up the system `PATH` if not found. You can verify its state in the logs:
- `Video transcoding enabled using ffmpeg: <path>` (Successful detection)
- `Video transcoding disabled (ffmpeg not found)` (Transcoding disabled gracefully)

## Project Structure

```
├── main.go              # Wails entry point
├── app.go               # Wails bindings (API bridge)
├── domain/              # Entity definitions (Media, Collection, PlayerConfig)
├── core/                # Business logic services
├── repository/          # SQLite data access layer
├── filesystem/          # File system operations (scanner, mover)
├── thumbs/              # 2-Tier thumbnail cache (LRU memory & persistent disk cache)
├── server/              # Standalone HTTP File Server with telemetry/metrics
├── transcode/           # Video HLS transcoding service
├── utils/               # Shared utilities (validation, formatting)
├── frontend/            # Svelte + Vite app
│   └── src/
│       ├── components/  # Reusable Svelte components
│       ├── pages/       # Page-level components
│       ├── stores/      # Svelte writable stores
│       └── services/    # Wails API wrappers
├── docs/                # Documentation
└── build/               # Build scripts (output built binary to build/bin/)
```

## Architecture

```
Frontend (Svelte) → Wails Bridge (app.go) → Core Services → Repository → SQLite
                                                           → File System
```

**Key principles:**
- File system is the source of truth for media
- SQLite stores metadata + indexing
- Services depend on repository interfaces
- Business logic is testable without UI

## How to Add a New Feature

1. **Define domain models** in `domain/`
2. **Add repository methods** in `repository/`
3. **Implement service logic** in `core/`
4. **Expose via Wails API** in `app.go`
5. **Create frontend page/component** in `frontend/src/`
6. **Write tests** alongside your service

## Running Tests

```bash
go test ./... -v -cover
```

## Database

SQLite database (`media_collector.db`) is auto-created on first launch.
Schema lives in `repository/schema.sql` and is embedded via `//go:embed`.
