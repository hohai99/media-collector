import { writable } from 'svelte/store';

export const currentPage = writable('home');
export const masterFolder = writable('');
export const isLoading = writable(false);
export const toasts = writable([]);

let toastId = 0;

export function showToast(message, type = 'info', duration = 3000) {
    const id = ++toastId;
    toasts.update(t => [...t, { id, message, type }]);
    setTimeout(() => {
        toasts.update(t => t.filter(x => x.id !== id));
    }, duration);
}

export function navigateTo(page) {
    currentPage.set(page);
}
