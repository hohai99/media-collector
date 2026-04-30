<script>
    import { onMount } from 'svelte';
    import CollectionCard from '../components/CollectionCard.svelte';
    import MediaGrid from '../components/MediaGrid.svelte';
    import Modal from '../components/Modal.svelte';
    import { selectedCollection, collections, buildTree } from '../stores/collections.js';
    import { mediaItems, selectedMedia, clearSelection, selectedCount } from '../stores/media.js';
    import { showToast, isLoading, masterFolder } from '../stores/app.js';
    import * as api from '../services/api.js';

    let showMoveModal = false;
    let moveTargetId = '';
    let topCollections = [];
    let subCollections = [];
    let breadcrumbs = [];
    let synced = false;
    let viewMode = 'gallery'; // 'gallery' or 'detail'

    onMount(async () => {
        if ($masterFolder) {
            await syncAndLoad();
        }
    });

    async function syncAndLoad() {
        isLoading.set(true);
        try {
            await api.syncAndScan();
            synced = true;

            // Refresh collections store
            const tree = await api.getCollectionTree();
            collections.set(tree || []);

            if ($selectedCollection) {
                await enterCollection($selectedCollection);
            } else {
                await loadTopLevel();
            }
        } catch (e) {
            console.error('Sync error:', e);
            showToast('Sync failed: ' + e, 'error');
        }
        isLoading.set(false);
    }

    async function loadTopLevel() {
        try {
            topCollections = await api.getTopLevelCollections() || [];
            topCollections.sort((a, b) => a.name.localeCompare(b.name));
            viewMode = 'gallery';
            breadcrumbs = [];
        } catch (e) {
            showToast('Failed to load collections: ' + e, 'error');
        }
    }

    async function enterCollection(col) {
        isLoading.set(true);
        selectedCollection.set(col);
        viewMode = 'detail';

        try {
            // Build breadcrumbs
            breadcrumbs = await buildBreadcrumbs(col);

            // Load sub-collections
            subCollections = await api.getSubCollections(col.id) || [];
            subCollections.sort((a, b) => a.name.localeCompare(b.name));

            // Load media in this collection (direct only)
            const items = await api.getMediaByCollection(col.id) || [];
            items.sort((a, b) => a.name.localeCompare(b.name));
            mediaItems.set(items);
        } catch (e) {
            showToast('Failed to load collection: ' + e, 'error');
        }
        isLoading.set(false);
    }

    async function buildBreadcrumbs(col) {
        const crumbs = [col];
        let current = col;
        // Walk up the parent chain
        while (current.parentId) {
            const allCols = $collections;
            const parent = allCols.find(c => c.id === current.parentId);
            if (parent) {
                crumbs.unshift(parent);
                current = parent;
            } else {
                break;
            }
        }
        return crumbs;
    }

    function goToGallery() {
        selectedCollection.set(null);
        clearSelection();
        loadTopLevel();
    }

    function handleCollectionClick(col) {
        clearSelection();
        enterCollection(col);
    }

    async function handleMoveMedia() {
        if (!moveTargetId || $selectedMedia.length === 0) return;
        try {
            await api.moveMedia($selectedMedia, moveTargetId);
            showToast(`Moved ${$selectedMedia.length} files`, 'success');
            clearSelection();
            showMoveModal = false;
            if ($selectedCollection) {
                await enterCollection($selectedCollection);
            }
        } catch (e) {
            showToast('Move failed: ' + e, 'error');
        }
    }

    // When selectedCollection changes from sidebar
    $: if ($selectedCollection && synced) {
        enterCollection($selectedCollection);
    }
</script>

<div class="page-container">
    {#if !$masterFolder}
        <div class="empty-state">
            <span class="icon">📂</span>
            <h3>No Master Folder Selected</h3>
            <p class="text-sm text-muted">Go to Import to select your master folder</p>
        </div>
    {:else if $isLoading}
        <div class="empty-state">
            <div class="spinner"></div>
            <p class="text-muted">Loading…</p>
        </div>
    {:else if viewMode === 'gallery'}
        <!-- All Collections Gallery -->
        <div class="page-header">
            <h1>📁 All Collections</h1>
            <span class="text-sm text-muted">{topCollections.length} collection{topCollections.length !== 1 ? 's' : ''}</span>
        </div>

        {#if topCollections.length === 0}
            <div class="empty-state">
                <span class="icon">📂</span>
                <h3>No Collections Found</h3>
                <p class="text-sm text-muted">Add folders to your master folder or use Import → Scan to discover media</p>
            </div>
        {:else}
            <div class="collection-grid">
                {#each topCollections as col (col.id)}
                    <CollectionCard collection={col} onClick={handleCollectionClick} />
                {/each}
            </div>
        {/if}

    {:else}
        <!-- Collection Detail View -->
        <div class="page-header">
            <div class="breadcrumb-row">
                <button class="btn btn-ghost btn-sm" on:click={goToGallery}>
                    📁 All Collections
                </button>
                {#each breadcrumbs as crumb, i}
                    <span class="breadcrumb-sep">›</span>
                    {#if i < breadcrumbs.length - 1}
                        <button class="btn btn-ghost btn-sm" on:click={() => handleCollectionClick(crumb)}>
                            {crumb.name}
                        </button>
                    {:else}
                        <span class="breadcrumb-current">{crumb.name}</span>
                    {/if}
                {/each}
            </div>

            <div class="flex gap-sm items-center">
                {#if $selectedCount > 0}
                    <span class="badge badge-image">{$selectedCount} selected</span>
                    <button class="btn btn-secondary btn-sm" on:click={() => showMoveModal = true}>
                        📦 Move
                    </button>
                    <button class="btn btn-ghost btn-sm" on:click={clearSelection}>
                        Clear
                    </button>
                {/if}
            </div>
        </div>

        <!-- Sub-collections -->
        {#if subCollections.length > 0}
            <div class="sub-section">
                <h3 class="section-label">📂 Sub-folders ({subCollections.length})</h3>
                <div class="collection-grid collection-grid-sm">
                    {#each subCollections as sub (sub.id)}
                        <CollectionCard collection={sub} onClick={handleCollectionClick} />
                    {/each}
                </div>
            </div>
        {/if}

        <!-- Media files -->
        <div class="sub-section">
            {#if subCollections.length > 0}
                <h3 class="section-label">🖼️ Media Files ({$mediaItems.length})</h3>
            {/if}
            <MediaGrid items={$mediaItems} />
        </div>
    {/if}
</div>

<!-- Move Modal -->
<Modal title="Move Media" show={showMoveModal} on:close={() => showMoveModal = false}>
    <p class="text-sm text-muted" style="margin-bottom: var(--space-md);">
        Select target collection for {$selectedCount} file(s):
    </p>
    <div class="move-target-list">
        {#each $collections as col}
            <label class="move-target-item" class:active={moveTargetId === col.id}>
                <input type="radio" bind:group={moveTargetId} value={col.id} />
                <span>📂 {col.name}</span>
            </label>
        {/each}
    </div>
    <div class="modal-actions">
        <button class="btn btn-secondary" on:click={() => showMoveModal = false}>Cancel</button>
        <button class="btn btn-primary" on:click={handleMoveMedia} disabled={!moveTargetId}>
            Move Files
        </button>
    </div>
</Modal>

<style>
    .collection-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
        gap: var(--space-lg);
    }
    .collection-grid-sm {
        grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
        gap: var(--space-md);
    }

    .breadcrumb-row {
        display: flex;
        align-items: center;
        gap: 2px;
        flex-wrap: wrap;
    }
    .breadcrumb-sep {
        color: var(--text-muted);
        font-size: 0.9rem;
        margin: 0 2px;
    }
    .breadcrumb-current {
        font-weight: 600;
        font-size: 0.875rem;
        color: var(--text-primary);
        padding: 4px 10px;
    }

    .sub-section {
        margin-bottom: var(--space-xl);
    }
    .section-label {
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--text-secondary);
        margin-bottom: var(--space-md);
    }

    .move-target-list {
        display: flex;
        flex-direction: column;
        gap: 4px;
        max-height: 300px;
        overflow-y: auto;
    }
    .move-target-item {
        display: flex;
        align-items: center;
        gap: var(--space-sm);
        padding: var(--space-sm);
        border-radius: var(--radius-sm);
        cursor: pointer;
        font-size: 0.875rem;
        color: var(--text-secondary);
        transition: background var(--duration-fast) var(--ease-out);
    }
    .move-target-item:hover {
        background: var(--bg-hover);
    }
    .move-target-item.active {
        background: var(--accent-muted);
        color: var(--text-accent);
    }
    .move-target-item input[type="radio"] {
        accent-color: var(--accent);
    }
</style>
