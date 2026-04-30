import { writable, derived } from 'svelte/store';

export const mediaItems = writable([]);
export const selectedMedia = writable([]);

export const selectedCount = derived(selectedMedia, $s => $s.length);

export function toggleMediaSelection(mediaId) {
    selectedMedia.update(items => {
        if (items.includes(mediaId)) {
            return items.filter(id => id !== mediaId);
        }
        return [...items, mediaId];
    });
}

export function clearSelection() {
    selectedMedia.set([]);
}

export function selectAll(ids) {
    selectedMedia.set([...ids]);
}
