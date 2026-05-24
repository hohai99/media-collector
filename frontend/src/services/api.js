/**
 * API service — thin wrappers around Wails-generated Go bindings.
 * Import paths are resolved to frontend/wailsjs/go/main/App at build time.
 */
import {
    SelectMasterFolder,
    GetMasterFolder,
    ScanMedia,
    GetAllMedia,
    GetMediaByCollection,
    GetMediaByCollectionRecursive,
    GetUncategorizedMedia,
    MoveMedia,
    CreateCollection,
    GetCollectionTree,
    GetTopLevelCollections,
    GetSubCollections,
    DeleteCollection,
    SyncCollections,
    SyncAndScan,
    GetCollectionThumbnail,
    GetCollectionMediaCount,
    CreatePlayerConfig,
    GetPlayerConfigs,
    GetPlayerConfig,
    DeletePlayerConfig,
    ResolvePlayerMedia,
    SelectSoundSource,
    GetAbsolutePath,
    GetFileServerURL,
    GetStreamBaseURL,
    HasTranscoder,
} from '../../wailsjs/go/main/App';

// ─── Folder ────────────────────────────────────────────────
export async function selectMasterFolder() {
    return await SelectMasterFolder();
}

export async function getMasterFolder() {
    return await GetMasterFolder();
}

// ─── Media ─────────────────────────────────────────────────
export async function scanMedia() {
    return await ScanMedia();
}

export async function getAllMedia() {
    return await GetAllMedia();
}

export async function getMediaByCollection(collectionId) {
    return await GetMediaByCollection(collectionId);
}

export async function getMediaByCollectionRecursive(collectionId) {
    return await GetMediaByCollectionRecursive(collectionId);
}

export async function getUncategorizedMedia() {
    return await GetUncategorizedMedia();
}

export async function moveMedia(mediaIds, collectionId) {
    return await MoveMedia(mediaIds, collectionId);
}

// ─── Collections ───────────────────────────────────────────
export async function createCollection(name, parentId = null) {
    return await CreateCollection(name, parentId || null);
}

export async function getCollectionTree() {
    return await GetCollectionTree();
}

export async function getTopLevelCollections() {
    return await GetTopLevelCollections();
}

export async function getSubCollections(parentId) {
    return await GetSubCollections(parentId);
}

export async function deleteCollection(id) {
    return await DeleteCollection(id);
}

export async function syncCollections() {
    return await SyncCollections();
}

export async function syncAndScan() {
    return await SyncAndScan();
}

export async function getCollectionThumbnail(collectionId) {
    return await GetCollectionThumbnail(collectionId);
}

export async function getCollectionMediaCount(collectionId) {
    return await GetCollectionMediaCount(collectionId);
}

// ─── Player ────────────────────────────────────────────────
export async function createPlayerConfig(input) {
    return await CreatePlayerConfig(input);
}

export async function getPlayerConfigs() {
    return await GetPlayerConfigs();
}

export async function getPlayerConfig(id) {
    return await GetPlayerConfig(id);
}

export async function deletePlayerConfig(id) {
    return await DeletePlayerConfig(id);
}

export async function resolvePlayerMedia(configId) {
    return await ResolvePlayerMedia(configId);
}

export async function selectSoundSource() {
    return await SelectSoundSource();
}

// ─── Utility ───────────────────────────────────────────────
export async function getAbsolutePath(path) {
    return await GetAbsolutePath(path);
}

export async function getFileServerURL() {
    return await GetFileServerURL();
}

export async function getStreamBaseURL() {
    return await GetStreamBaseURL();
}

export async function hasTranscoder() {
    return await HasTranscoder();
}

