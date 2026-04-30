import { writable } from 'svelte/store';

export const playerConfigs = writable([]);
export const activeConfig = writable(null);
export const playerMedia = writable([]);
export const isPlaying = writable(false);
export const currentIndex = writable(0);
