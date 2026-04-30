<script>
    import { showToast, isLoading, masterFolder } from '../stores/app.js';
    import { collections } from '../stores/collections.js';
    import { mediaItems } from '../stores/media.js';
    import * as api from '../services/api.js';
    import MediaGrid from '../components/MediaGrid.svelte';

    let scanResults = [];
    let hasScanned = false;

    async function handleSelectFolder() {
        try {
            const folder = await api.selectMasterFolder();
            if (folder) {
                masterFolder.set(folder);
                showToast('Master folder set: ' + folder, 'success');
                // Refresh collections
                const tree = await api.getCollectionTree();
                collections.set(tree || []);
            }
        } catch (e) {
            showToast('Failed to select folder: ' + e, 'error');
        }
    }

    async function handleScan() {
        if (!$masterFolder) {
            showToast('Please select a master folder first', 'error');
            return;
        }
        isLoading.set(true);
        hasScanned = false;
        try {
            scanResults = await api.scanMedia() || [];
            hasScanned = true;
            showToast(`Found ${scanResults.length} media files`, 'success');

            // Refresh global media list
            const allMedia = await api.getAllMedia();
            mediaItems.set(allMedia || []);

            // Refresh collections (sync may have found new folders)
            const tree = await api.getCollectionTree();
            collections.set(tree || []);
        } catch (e) {
            showToast('Scan failed: ' + e, 'error');
        }
        isLoading.set(false);
    }

    let filterType = 'all';
    $: filteredResults = filterType === 'all'
        ? scanResults
        : scanResults.filter(m => m.type === filterType);

    $: imageCount = scanResults.filter(m => m.type === 'image').length;
    $: videoCount = scanResults.filter(m => m.type === 'video').length;
    $: audioCount = scanResults.filter(m => m.type === 'audio').length;
</script>

<div class="page-container">
    <div class="page-header">
        <h1>📥 Import Media</h1>
        <div class="flex gap-sm">
            <button class="btn btn-secondary" on:click={handleSelectFolder}>
                📁 {$masterFolder ? 'Change Folder' : 'Select Master Folder'}
            </button>
            <button class="btn btn-primary" on:click={handleScan} disabled={!$masterFolder}>
                🔍 Scan Media
            </button>
        </div>
    </div>

    {#if !$masterFolder}
        <div class="empty-state">
            <span class="icon">📂</span>
            <h3>No Master Folder Selected</h3>
            <p class="text-sm text-muted">Select a master folder to begin scanning for media files.</p>
            <button class="btn btn-primary btn-lg" on:click={handleSelectFolder}>
                📁 Select Master Folder
            </button>
        </div>
    {:else if $isLoading}
        <div class="empty-state">
            <div class="spinner"></div>
            <h3>Scanning…</h3>
            <p class="text-muted">Looking for media files in your master folder</p>
        </div>
    {:else if hasScanned}
        <div class="scan-summary glass-card" style="margin-bottom: var(--space-lg);">
            <div class="summary-grid">
                <div class="summary-item">
                    <span class="summary-number">{scanResults.length}</span>
                    <span class="summary-label text-muted">Total Files</span>
                </div>
                <div class="summary-item">
                    <span class="summary-number" style="color: #a29bfe;">{imageCount}</span>
                    <span class="summary-label text-muted">Images</span>
                </div>
                <div class="summary-item">
                    <span class="summary-number" style="color: #00cec9;">{videoCount}</span>
                    <span class="summary-label text-muted">Videos</span>
                </div>
                <div class="summary-item">
                    <span class="summary-number" style="color: #fdcb6e;">{audioCount}</span>
                    <span class="summary-label text-muted">Audio</span>
                </div>
            </div>
        </div>

        <!-- Type filter -->
        <div class="filter-bar" style="margin-bottom: var(--space-md);">
            <button class="btn" class:btn-primary={filterType === 'all'} class:btn-ghost={filterType !== 'all'} on:click={() => filterType = 'all'}>All</button>
            <button class="btn" class:btn-primary={filterType === 'image'} class:btn-ghost={filterType !== 'image'} on:click={() => filterType = 'image'}>🖼️ Images</button>
            <button class="btn" class:btn-primary={filterType === 'video'} class:btn-ghost={filterType !== 'video'} on:click={() => filterType = 'video'}>🎬 Videos</button>
            <button class="btn" class:btn-primary={filterType === 'audio'} class:btn-ghost={filterType !== 'audio'} on:click={() => filterType = 'audio'}>🎵 Audio</button>
        </div>

        <MediaGrid items={filteredResults} />
    {:else}
        <div class="empty-state">
            <span class="icon">🔍</span>
            <h3>Ready to Scan</h3>
            <p class="text-sm text-muted">Click "Scan Media" to discover files in your master folder</p>
        </div>
    {/if}
</div>

<style>
    .scan-summary {
        padding: var(--space-lg);
    }

    .summary-grid {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: var(--space-lg);
        text-align: center;
    }

    .summary-item {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .summary-number {
        font-size: 2rem;
        font-weight: 700;
        color: var(--text-primary);
    }

    .summary-label {
        font-size: 0.8rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .filter-bar {
        display: flex;
        gap: var(--space-xs);
    }
</style>
