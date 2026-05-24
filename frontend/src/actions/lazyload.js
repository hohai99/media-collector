/**
 * Svelte action: lazy-load images via IntersectionObserver.
 *
 * Usage:
 *   <img use:lazyload={thumbUrl} alt="..." class="..." />
 *
 * The action keeps src empty until the element enters the viewport
 * (with 300px of margin). When the src prop changes, the old load is
 * cancelled and a new observation begins.
 */
export function lazyload(node, src) {
    let currentSrc = src;
    let loaded = false;

    const observer = new IntersectionObserver(
        ([entry]) => {
            if (entry.isIntersecting && !loaded) {
                loaded = true;
                observer.unobserve(node);
                node.src = currentSrc;
            }
        },
        { rootMargin: '300px' }
    );

    if (currentSrc) {
        observer.observe(node);
    }

    return {
        update(newSrc) {
            if (newSrc !== currentSrc) {
                currentSrc = newSrc;
                loaded = false;
                node.removeAttribute('src');
                if (currentSrc) {
                    observer.observe(node);
                }
            }
        },
        destroy() {
            observer.disconnect();
        },
    };
}
