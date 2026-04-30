package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"media-collector/core"
	"media-collector/domain"
	"media-collector/filesystem"
	"media-collector/repository"
	"media-collector/utils"
)

// App struct holds all services and is the Wails binding target.
type App struct {
	ctx context.Context
	db  *sql.DB

	// Services
	mediaSvc      *core.MediaService
	collectionSvc *core.CollectionService
	playerSvc     *core.PlayerService

	// Repositories (exposed for direct config access)
	configRepo *repository.ConfigRepository
	collRepo   *repository.CollectionRepository
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. It initialises the database and
// all service layers.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Database path next to the executable (or in dev, current dir)
	dbPath := "media_collector.db"
	db, err := repository.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	a.db = db

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Repositories
	mediaRepo := repository.NewMediaRepository(db)
	collRepo := repository.NewCollectionRepository(db)
	configRepo := repository.NewConfigRepository(db)
	playerRepo := repository.NewPlayerRepository(db)

	a.configRepo = configRepo
	a.collRepo = collRepo

	// File system
	scanner := filesystem.NewScanner()
	mover := filesystem.NewMover()

	// Services
	a.mediaSvc = core.NewMediaService(mediaRepo, scanner, mover, collRepo)
	a.collectionSvc = core.NewCollectionService(collRepo)
	a.playerSvc = core.NewPlayerService(playerRepo, mediaRepo, collRepo)
}

// beforeClose is called when the application is about to quit.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.db != nil {
		a.db.Close()
	}
	return false
}

// ---------------------------------------------------------------------------
// Folder Management
// ---------------------------------------------------------------------------

// SelectMasterFolder opens a native folder dialog and stores the selection.
func (a *App) SelectMasterFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Master Folder",
	})
	if err != nil {
		return "", fmt.Errorf("open folder dialog: %w", err)
	}
	if folder == "" {
		return "", nil // user cancelled
	}

	if err := a.configRepo.Set(domain.ConfigKeyMasterFolder, folder); err != nil {
		return "", fmt.Errorf("save master folder: %w", err)
	}

	// Sync collections from the selected folder
	if err := a.collectionSvc.SyncCollections(folder); err != nil {
		return "", fmt.Errorf("sync collections: %w", err)
	}

	return folder, nil
}

// GetMasterFolder returns the currently configured master folder path.
func (a *App) GetMasterFolder() (string, error) {
	return a.configRepo.Get(domain.ConfigKeyMasterFolder)
}

// ---------------------------------------------------------------------------
// Media
// ---------------------------------------------------------------------------

// ScanMedia scans the master folder for media files and imports them into the DB.
func (a *App) ScanMedia() ([]domain.Media, error) {
	masterFolder, err := a.configRepo.Get(domain.ConfigKeyMasterFolder)
	if err != nil || masterFolder == "" {
		return nil, fmt.Errorf("master folder not set")
	}
	return a.mediaSvc.ScanFolder(masterFolder)
}

// GetAllMedia returns every media item in the database.
func (a *App) GetAllMedia() ([]domain.Media, error) {
	return a.mediaSvc.GetAllMedia()
}

// GetMediaByCollection returns media belonging to a specific collection.
func (a *App) GetMediaByCollection(collectionID string) ([]domain.Media, error) {
	return a.mediaSvc.GetMediaByCollection(collectionID)
}

// GetMediaByCollectionRecursive returns media from a collection and all its sub-collections.
func (a *App) GetMediaByCollectionRecursive(collectionID string) ([]domain.Media, error) {
	return a.mediaSvc.GetMediaByCollectionRecursive(collectionID)
}

// GetUncategorizedMedia returns media not in any collection.
func (a *App) GetUncategorizedMedia() ([]domain.Media, error) {
	return a.mediaSvc.GetUncategorized()
}

// MoveMedia moves the specified media files into a target collection.
func (a *App) MoveMedia(mediaIDs []string, collectionID string) error {
	coll, err := a.collRepo.GetByID(collectionID)
	if err != nil {
		return fmt.Errorf("get collection: %w", err)
	}
	if coll == nil {
		return fmt.Errorf("collection %q not found", collectionID)
	}
	return a.mediaSvc.MoveMedia(mediaIDs, coll.Path, coll.ID)
}

// ---------------------------------------------------------------------------
// Collections
// ---------------------------------------------------------------------------

// CreateCollection creates a new collection (folder) in the master folder.
func (a *App) CreateCollection(name string, parentID *string) (*domain.Collection, error) {
	if err := utils.ValidateName(name); err != nil {
		return nil, err
	}
	masterFolder, err := a.configRepo.Get(domain.ConfigKeyMasterFolder)
	if err != nil || masterFolder == "" {
		return nil, fmt.Errorf("master folder not set")
	}
	return a.collectionSvc.CreateCollection(name, parentID, masterFolder)
}

// GetCollectionTree returns all collections as a flat list.
func (a *App) GetCollectionTree() ([]domain.Collection, error) {
	return a.collectionSvc.GetCollectionTree()
}

// GetTopLevelCollections returns only root-level collections (no parent).
func (a *App) GetTopLevelCollections() ([]domain.Collection, error) {
	return a.collRepo.GetTopLevel()
}

// GetSubCollections returns direct child collections of a parent.
func (a *App) GetSubCollections(parentID string) ([]domain.Collection, error) {
	return a.collRepo.GetChildren(parentID)
}

// DeleteCollection removes a collection from the DB (folder stays on disk).
func (a *App) DeleteCollection(id string) error {
	return a.collectionSvc.DeleteCollection(id)
}

// SyncCollections reconciles filesystem folders with the DB.
func (a *App) SyncCollections() error {
	masterFolder, err := a.configRepo.Get(domain.ConfigKeyMasterFolder)
	if err != nil || masterFolder == "" {
		return fmt.Errorf("master folder not set")
	}
	return a.collectionSvc.SyncCollections(masterFolder)
}

// SyncAndScan syncs collections from disk then scans for media files.
// This detects manually copied folders and their media in one call.
func (a *App) SyncAndScan() error {
	masterFolder, err := a.configRepo.Get(domain.ConfigKeyMasterFolder)
	if err != nil || masterFolder == "" {
		return fmt.Errorf("master folder not set")
	}

	if err := a.collectionSvc.SyncCollections(masterFolder); err != nil {
		return fmt.Errorf("sync collections: %w", err)
	}

	if _, err := a.mediaSvc.ScanFolder(masterFolder); err != nil {
		return fmt.Errorf("scan media: %w", err)
	}

	return nil
}

// GetCollectionThumbnail returns the path of the first image in a collection.
func (a *App) GetCollectionThumbnail(collectionID string) (string, error) {
	coll, err := a.collRepo.GetByID(collectionID)
	if err != nil {
		return "", fmt.Errorf("get collection: %w", err)
	}
	if coll == nil {
		return "", fmt.Errorf("collection %q not found", collectionID)
	}
	return a.collectionSvc.GetThumbnail(coll.Path), nil
}

// GetCollectionMediaCount returns the number of media files in a collection (direct only).
func (a *App) GetCollectionMediaCount(collectionID string) (int, error) {
	media, err := a.mediaSvc.GetMediaByCollection(collectionID)
	if err != nil {
		return 0, err
	}
	return len(media), nil
}

// ---------------------------------------------------------------------------
// Player
// ---------------------------------------------------------------------------

// CreatePlayerConfig creates a new player configuration.
func (a *App) CreatePlayerConfig(input domain.CreateConfigInput) (*domain.PlayerConfig, error) {
	if err := utils.ValidateName(input.Name); err != nil {
		return nil, err
	}
	return a.playerSvc.CreateConfig(input)
}

// GetPlayerConfigs returns all player configurations.
func (a *App) GetPlayerConfigs() ([]domain.PlayerConfig, error) {
	return a.playerSvc.GetAllConfigs()
}

// GetPlayerConfig returns a single player config by ID.
func (a *App) GetPlayerConfig(id string) (*domain.PlayerConfig, error) {
	return a.playerSvc.GetConfig(id)
}

// DeletePlayerConfig removes a player configuration.
func (a *App) DeletePlayerConfig(id string) error {
	return a.playerSvc.DeleteConfig(id)
}

// ResolvePlayerMedia resolves the full media list for a player config.
func (a *App) ResolvePlayerMedia(configID string) ([]domain.Media, error) {
	return a.playerSvc.ResolveMediaList(configID)
}

// ---------------------------------------------------------------------------
// Utility — exposed to frontend for file path operations
// ---------------------------------------------------------------------------

// GetAbsolutePath converts a relative path to absolute (used for media src).
func (a *App) GetAbsolutePath(path string) string {
	abs, err := filepath.Abs(filepath.FromSlash(path))
	if err != nil {
		return path
	}
	return filepath.ToSlash(abs)
}

