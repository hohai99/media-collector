<script>
    import { onMount } from 'svelte';
    import { showToast, navigateTo } from '../stores/app.js';
    import { playerConfigs, activeConfig, playerMedia } from '../stores/player.js';
    import { collections } from '../stores/collections.js';
    import FolderTree from '../components/FolderTree.svelte';
    import * as api from '../services/api.js';
    import { formatDuration } from '../services/utils.js';

    let configName = '';
    let timePerPicture = 5;
    let soundSource = '';
    let selectedCollections = [];
    let configs = [];

    // For estimating image count
    let estimatedImageCount = 0;
    let loadingEstimate = false;

    onMount(async () => {
        await loadConfigs();
    });

    async function loadConfigs() {
        try {
            configs = await api.getPlayerConfigs() || [];
            playerConfigs.set(configs);
        } catch (e) {
            showToast('Failed to load configs: ' + e, 'error');
        }
    }

    // When collections change, estimate the image count
    async function updateEstimate() {
        if (selectedCollections.length === 0) {
            estimatedImageCount = 0;
            return;
        }
        loadingEstimate = true;
        try {
            let count = 0;
            for (const colId of selectedCollections) {
                const media = await api.getMediaByCollectionRecursive(colId);
                if (media) {
                    count += media.filter(m => m.type === 'image').length;
                }
            }
            estimatedImageCount = count;
        } catch (e) {
            estimatedImageCount = 0;
        }
        loadingEstimate = false;
    }

    async function handleCreate() {
        if (!configName.trim()) {
            showToast('Please enter a config name', 'error');
            return;
        }
        if (selectedCollections.length === 0) {
            showToast('Select at least one collection', 'error');
            return;
        }

        try {
            await api.createPlayerConfig({
                name: configName.trim(),
                collections: selectedCollections,
                mediaIds: [],
                timePerPicture: timePerPicture,
                soundSource: soundSource,
            });
            showToast('Player config created!', 'success');
            configName = '';
            selectedCollections = [];
            timePerPicture = 5;
            soundSource = '';
            estimatedImageCount = 0;
            await loadConfigs();
        } catch (e) {
            showToast('Create failed: ' + e, 'error');
        }
    }

    async function handlePlay(config) {
        try {
            const media = await api.resolvePlayerMedia(config.id);
            if (!media || media.length === 0) {
                showToast('No media found for this config', 'error');
                return;
            }
            activeConfig.set(config);
            playerMedia.set(media);
            navigateTo('player');
        } catch (e) {
            showToast('Failed to load media: ' + e, 'error');
        }
    }

    async function handleDelete(id) {
        try {
            await api.deletePlayerConfig(id);
            showToast('Config deleted', 'success');
            await loadConfigs();
        } catch (e) {
            showToast('Delete failed: ' + e, 'error');
        }
    }

    function handleCollectionChange(e) {
        selectedCollections = e.detail;
        updateEstimate();
    }

    async function handleSelectSound() {
        try {
            const file = await api.selectSoundSource();
            if (file) {
                soundSource = file;
            }
        } catch (e) {
            showToast('Failed to select sound: ' + e, 'error');
        }
    }

    function clearSound() {
        soundSource = '';
    }

    $: estimatedTotal = estimatedImageCount * timePerPicture;
</script>

<div class="page-container">
    <div class="page-header">
        <h1>🎛️ Player Configurations</h1>
    </div>

    <!-- Create Form -->
    <div class="glass-card" style="margin-bottom: var(--space-xl);">
        <h3 style="margin-bottom: var(--space-md);">Create New Config</h3>

        <div class="form-grid">
            <div class="form-group">
                <label class="label" for="config-name">Config Name</label>
                <input
                    id="config-name"
                    class="input"
                    bind:value={configName}
                    placeholder="e.g. Morning Slideshow"
                />
            </div>

            <div class="form-group">
                <label class="label" for="time-per-picture">Time Per Picture (seconds)</label>
                <input
                    id="time-per-picture"
                    class="input"
                    type="number"
                    min="1"
                    max="999"
                    bind:value={timePerPicture}
                />
            </div>
        </div>

        <!-- Sound Source -->
        <div class="form-group" style="margin-top: var(--space-md);">
            <label class="label">Sound Source (optional)</label>
            <div class="sound-picker">
                <button class="btn btn-secondary btn-sm" on:click={handleSelectSound}>
                    🎵 {soundSource ? 'Change Sound' : 'Select Sound File'}
                </button>
                {#if soundSource}
                    <div class="sound-info">
                        <span class="text-sm truncate" title={soundSource}>
                            📎 {soundSource.split(/[\\/]/).pop()}
                        </span>
                        <button class="btn btn-ghost btn-sm" on:click={clearSound} title="Remove">✕</button>
                    </div>
                {/if}
            </div>
            <span class="text-sm text-muted" style="margin-top: 4px;">
                Plays as background music. Loops if shorter than slideshow, stops at end if longer.
            </span>
        </div>

        <div style="margin-top: var(--space-md);">
            <label class="label">Select Collections</label>
            <FolderTree selectedIds={selectedCollections} on:change={handleCollectionChange} />
        </div>

        {#if selectedCollections.length > 0}
            <div class="estimate-bar" style="margin-top: var(--space-md);">
                {#if loadingEstimate}
                    <span class="text-sm text-muted">Counting media…</span>
                {:else}
                    <span class="text-sm text-muted">
                        🖼️ {estimatedImageCount} image{estimatedImageCount !== 1 ? 's' : ''}
                        × {timePerPicture}s
                    </span>
                    <span class="text-sm" style="font-weight: 600; color: var(--text-accent);">
                        = {formatDuration(estimatedTotal)} total
                    </span>
                {/if}
            </div>
        {/if}

        <div style="margin-top: var(--space-lg);">
            <button class="btn btn-primary" on:click={handleCreate}>
                ✨ Create Config
            </button>
        </div>
    </div>

    <!-- Existing Configs -->
    {#if configs.length > 0}
        <h3 style="margin-bottom: var(--space-md);">Saved Configurations</h3>
        <div class="config-list">
            {#each configs as cfg}
                <div class="config-card glass-card">
                    <div class="config-header">
                        <div>
                            <h4>{cfg.name}</h4>
                            <span class="text-sm text-muted">
                                {cfg.timePerPicture || cfg.transitionTime}s/slide
                                · {formatDuration(cfg.totalPlayTime)} total
                                {#if cfg.soundSource}
                                    · 🎵 {cfg.soundSource.split(/[\\/]/).pop()}
                                {/if}
                            </span>
                        </div>
                        <div class="flex gap-sm">
                            <button class="btn btn-primary btn-sm" on:click={() => handlePlay(cfg)}>
                                ▶️ Play
                            </button>
                            <button class="btn btn-danger btn-sm" on:click={() => handleDelete(cfg.id)}>
                                🗑
                            </button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {:else}
        <div class="empty-state" style="padding: var(--space-xl) 0;">
            <span class="icon">🎬</span>
            <h3>No Configurations Yet</h3>
            <p class="text-sm text-muted">Create your first player config above</p>
        </div>
    {/if}
</div>

<style>
    .form-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--space-md);
    }

    .form-group {
        display: flex;
        flex-direction: column;
    }

    .sound-picker {
        display: flex;
        align-items: center;
        gap: var(--space-sm);
        flex-wrap: wrap;
    }

    .sound-info {
        display: flex;
        align-items: center;
        gap: var(--space-xs);
        background: var(--bg-tertiary);
        border: 1px solid var(--glass-border);
        border-radius: var(--radius-sm);
        padding: 4px 8px 4px 12px;
        max-width: 300px;
    }

    .estimate-bar {
        display: flex;
        align-items: center;
        gap: var(--space-md);
        padding: var(--space-sm) var(--space-md);
        background: var(--bg-tertiary);
        border-radius: var(--radius-sm);
        border: 1px solid var(--glass-border);
    }

    .config-list {
        display: flex;
        flex-direction: column;
        gap: var(--space-sm);
    }

    .config-card {
        padding: var(--space-md) var(--space-lg);
    }

    .config-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
</style>
