import { writable } from 'svelte/store';

export const collections = writable([]);
export const selectedCollection = writable(null);

/**
 * Build a tree structure from the flat collections list.
 * Each node gets a `children` array.
 */
export function buildTree(flatList) {
    const map = {};
    const roots = [];

    flatList.forEach(c => {
        map[c.id] = { ...c, children: [] };
    });

    flatList.forEach(c => {
        if (c.parentId && map[c.parentId]) {
            map[c.parentId].children.push(map[c.id]);
        } else {
            roots.push(map[c.id]);
        }
    });

    return roots;
}
