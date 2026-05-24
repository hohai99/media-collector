<script>
    import { onMount, onDestroy } from 'svelte';
    import MediaCard from './MediaCard.svelte';

    export let items = [];

    // Virtual scrolling config
    const VIRTUAL_THRESHOLD = 200;
    const CARD_HEIGHT = 220; // approximate height of a card in px
    const BUFFER = 10;       // extra rows to render above/below viewport

    let container;
    let scrollTop = 0;
    let containerHeight = 0;
    let columnsPerRow = 4;
    let useVirtual = false;

    $: useVirtual = items.length > VIRTUAL_THRESHOLD;
    $: totalRows = useVirtual ? Math.ceil(items.length / columnsPerRow) : 0;
    $: totalHeight = totalRows * CARD_HEIGHT;

    // Visible range calculation
    $: startRow = useVirtual ? Math.max(0, Math.floor(scrollTop / CARD_HEIGHT) - BUFFER) : 0;
    $: endRow = useVirtual
        ? Math.min(totalRows, Math.ceil((scrollTop + containerHeight) / CARD_HEIGHT) + BUFFER)
        : 0;
    $: startIndex = startRow * columnsPerRow;
    $: endIndex = useVirtual ? Math.min(items.length, endRow * columnsPerRow) : items.length;
    $: visibleItems = useVirtual ? items.slice(startIndex, endIndex) : items;
    $: topPad = startRow * CARD_HEIGHT;
    $: bottomPad = useVirtual ? Math.max(0, (totalRows - endRow) * CARD_HEIGHT) : 0;

    function handleScroll() {
        if (!container) return;
        scrollTop = container.scrollTop;
    }

    function updateColumns() {
        if (!container) return;
        containerHeight = container.clientHeight;
        const containerWidth = container.clientWidth;
        // Match CSS: minmax(200px, 1fr) with 16px gap
        const minCardWidth = 200;
        const gap = 16;
        columnsPerRow = Math.max(1, Math.floor((containerWidth + gap) / (minCardWidth + gap)));
    }

    let resizeObserver;

    onMount(() => {
        if (container) {
            updateColumns();
            resizeObserver = new ResizeObserver(() => updateColumns());
            resizeObserver.observe(container);
        }
    });

    onDestroy(() => {
        if (resizeObserver) resizeObserver.disconnect();
    });
</script>

{#if items.length === 0}
    <div class="empty-state">
        <span class="icon">📭</span>
        <h3>No media files</h3>
        <p class="text-sm text-muted">Scan your master folder to discover media files</p>
    </div>
{:else if useVirtual}
    <!-- Virtual scrolling for large lists -->
    <div class="virtual-container" bind:this={container} on:scroll={handleScroll}>
        <div class="virtual-spacer" style="height: {totalHeight}px;">
            <div class="media-grid" style="transform: translateY({topPad}px);">
                {#each visibleItems as item (item.id)}
                    <MediaCard media={item} />
                {/each}
            </div>
        </div>
    </div>
{:else}
    <div class="media-grid" bind:this={container}>
        {#each items as item (item.id)}
            <MediaCard media={item} />
        {/each}
    </div>
{/if}

<style>
    .virtual-container {
        height: 100%;
        overflow-y: auto;
        overflow-x: hidden;
    }

    .virtual-spacer {
        position: relative;
        width: 100%;
    }
</style>
