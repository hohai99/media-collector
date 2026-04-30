# Media Collector — Developer Guide

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | ≥ 1.21 | https://go.dev/dl/ |
| Node.js | ≥ 18 | https://nodejs.org/ |
| Wails CLI | ≥ 2.x | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |

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

## Project Structure

```
├── main.go              # Wails entry point
├── app.go               # Wails bindings (API bridge)
├── domain/              # Entity definitions (Media, Collection, PlayerConfig)
├── core/                # Business logic services
├── repository/          # SQLite data access layer
├── filesystem/          # File system operations (scanner, mover)
├── utils/               # Shared utilities (validation, formatting)
├── frontend/            # Svelte + Vite app
│   └── src/
│       ├── components/  # Reusable Svelte components
│       ├── pages/       # Page-level components
│       ├── stores/      # Svelte writable stores
│       └── services/    # Wails API wrappers
├── docs/                # Documentation
└── build/               # Build scripts
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
