# Media Collector

A desktop application for organizing, viewing, and playing your local media files (images, videos, and audio). Built with [Wails](https://wails.io/), it features a fast, lightweight Go backend and a modern Svelte frontend.

## Features

- **Local Media Management**: Keep all your files organized locally without needing cloud storage.
- **Master Folder Sync**: Select a "Master Folder" on your machine. The app automatically scans and syncs the directory structure and media files. Manually copied folders are detected automatically upon startup or refresh.
- **Collection Gallery**: Folders are represented as "Collections" with auto-generated cover thumbnails (using the first image in the folder).
- **Media Preview**: A built-in lightbox allows you to quickly view images, watch videos, and listen to audio directly in the app.
- **Custom Player**: Build custom slideshows or playlists by selecting multiple collections or specific media files. It calculates transition times and recursively gathers media from sub-folders.
- **Move Media**: Easily move files between collections directly from the app interface.

## Tech Stack

- **Backend**: Go
- **Frontend**: Svelte, JavaScript, HTML/CSS (Vite)
- **Framework**: Wails v2 (cross-platform desktop application framework)
- **Database**: SQLite (local `media_collector.db`)

## Prerequisites

- [Go](https://go.dev/doc/install) 1.20+
- [Node.js](https://nodejs.org/en/download/) 16+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

## Getting Started

### Live Development

To run the app in live development mode with hot-reloading:

```bash
wails dev
```

This will start a Vite development server for the frontend and compile the Go backend. Any changes to the Svelte components will automatically reflect in the application window.

### Building for Production

To build a standalone, redistributable executable for your OS:

```bash
wails build
```

The compiled application will be available in the `build/bin/` directory.

## Architecture Notes

- **Database**: The app creates a local `media_collector.db` SQLite database in the working directory to index files and relationships.
- **File Serving**: Because Wails WebView cannot natively load raw OS file paths (e.g., `C:/images/pic.jpg`) due to security restrictions, the app uses a custom Wails `AssetServer` HTTP Handler (`FileLoader` in `main.go`) to serve media to the frontend via `/localfile/` encoded URLs.
- **Recursive Collections**: The database supports a parent-child relationship for folders. Media resolution for the built-in Player recursively fetches media from all sub-folders of a selected collection.

## License

MIT License
