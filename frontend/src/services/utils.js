/**
 * Utility functions for the frontend.
 */

/**
 * Format file size in human-readable form.
 */
export function formatSize(bytes) {
    if (bytes === 0) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return (bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0) + ' ' + units[i];
}

/**
 * Format seconds as "Xm Ys" or "Xh Ym Zs".
 */
export function formatDuration(seconds) {
    if (!seconds || seconds <= 0) return '0s';
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    if (h > 0) return `${h}h ${m}m ${s}s`;
    if (m > 0) return `${m}m ${s}s`;
    return `${s}s`;
}

/**
 * Get the icon/emoji for a media type.
 */
export function mediaTypeIcon(type) {
    switch (type) {
        case 'image': return '🖼️';
        case 'video': return '🎬';
        case 'audio': return '🎵';
        default: return '📄';
    }
}

/**
 * Get the file extension from a path.
 */
export function getExtension(path) {
    const parts = path.split('.');
    return parts.length > 1 ? '.' + parts.pop().toLowerCase() : '';
}

/**
 * Debounce a function call.
 */
export function debounce(fn, delay = 300) {
    let timer;
    return (...args) => {
        clearTimeout(timer);
        timer = setTimeout(() => fn(...args), delay);
    };
}

/**
 * Cache for the file server base URL.
 * Fetched once from the Go backend on first use.
 */
let _fileServerBase = null;

/**
 * Initialize the file server URL cache.
 * Must be called once at app startup (e.g. in App.svelte onMount).
 */
export async function initFileServer() {
    const { getFileServerURL } = await import('../services/api.js');
    _fileServerBase = await getFileServerURL();
}

/**
 * Convert a raw filesystem path to a URL served by the standalone file server.
 * e.g. "F:/folder/img.jpg" → "http://127.0.0.1:12345/file/F%3A/folder/img.jpg"
 * Returns empty string if path is empty or server not initialized.
 */
export function localFileUrl(path) {
    if (!path) return '';
    if (!_fileServerBase) {
        // Fallback until initFileServer() completes
        const normalized = path.replace(/\\/g, '/');
        const encoded = normalized.split('/').map(s => encodeURIComponent(s)).join('/');
        return '/localfile/' + encoded;
    }
    // Normalise to forward slashes and encode each segment
    const normalized = path.replace(/\\/g, '/');
    const encoded = normalized.split('/').map(s => encodeURIComponent(s)).join('/');
    return _fileServerBase + encoded;
}

/**
 * Build a thumbnail URL served by the standalone file server.
 * e.g. "http://127.0.0.1:12345/thumb?path=F%3A%2Ffolder%2Fimg.jpg&w=300&h=300&fit=cover"
 * Returns empty string if path is empty or server not initialized.
 */
export function localThumbUrl(path, w = 300, h = 300, fit = 'cover') {
    if (!path || !_fileServerBase) return '';
    const thumbBase = _fileServerBase.replace('/file/', '/thumb');
    const normalized = path.replace(/\\/g, '/');
    return `${thumbBase}?path=${encodeURIComponent(normalized)}&w=${w}&h=${h}&fit=${fit}`;
}

/**
 * Build a stream URL for video transcoding (HLS).
 * e.g. "http://127.0.0.1:12345/stream?path=F%3A%2Ffolder%2Fvideo.mkv"
 */
export function localStreamUrl(path) {
    if (!path || !_fileServerBase) return '';
    const streamBase = _fileServerBase.replace('/file/', '/stream');
    const normalized = path.replace(/\\/g, '/');
    return `${streamBase}?path=${encodeURIComponent(normalized)}`;
}

/**
 * Check if a video file extension requires transcoding (not natively
 * supported by browsers). MP4 and WebM are natively supported.
 */
const NATIVE_VIDEO_EXTS = new Set(['.mp4', '.webm']);

export function needsTranscode(path) {
    if (!path) return false;
    const ext = getExtension(path);
    const videoExts = new Set(['.mp4', '.avi', '.mkv', '.mov', '.wmv', '.flv', '.webm']);
    if (!videoExts.has(ext)) return false;
    return !NATIVE_VIDEO_EXTS.has(ext);
}
