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
 * Convert a raw filesystem path to a URL the Wails WebView2 can load.
 * e.g. "F:/Media Collector/master-folder/img.jpg" → "/localfile/F:/Media Collector/master-folder/img.jpg"
 * Returns empty string if path is empty.
 */
export function localFileUrl(path) {
    if (!path) return '';
    // Normalise to forward slashes
    const normalized = path.replace(/\\/g, '/');
    // Encode each path segment to handle spaces, unicode, etc.
    // but keep the '/' separators intact
    const encoded = normalized.split('/').map(segment => encodeURIComponent(segment)).join('/');
    return '/localfile/' + encoded;
}

