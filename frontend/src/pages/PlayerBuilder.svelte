<script>
    import { onMount } from 'svelte';
    import { showToast, navigateTo } from '../stores/app.js';
    import { playerConfigs, activeConfig, playerMedia } from '../stores/player.js';
    import { collections } from '../stores/collections.js';
    import FolderTree from '../components/FolderTree.svelte';
    import * as api from '../services/api.js';
    import { formatDuration } from '../services/utils.js';

    let configName = '';
    let totalMinutes = 3;
    let selectedCollections = [];
    let configs = [];

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
                totalPlayTime: totalMinutes * 60,
            });
            showToast('Player config created!', 'success');
            configName = '';
            selectedCollections = [];
            totalMinutes = 3;
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
    }

    $: estimatedTransition = selectedCollections.length > 0
        ? Math.floor((totalMinutes * 60) / Math.max(selectedCollections.length, 1))
        : 0;
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
                <label class="label" for="total-time">Total Play Time (minutes)</label>
                <input
                    id="total-time"
                    class="input"
                    type="number"
                    min="1"
                    max="999"
                    bind:value={totalMinutes}
                />
                <span class="text-sm text-muted" style="margin-top: 4px;">
                    = {formatDuration(totalMinutes * 60)}
                </span>
            </div>
        </div>

        <div style="margin-top: var(--space-md);">
            <label class="label">Select Collections</label>
            <FolderTree selectedIds={selectedCollections} on:change={handleCollectionChange} />
        </div>

        {#if selectedCollections.length > 0}
            <div class="text-sm text-muted" style="margin-top: var(--space-sm);">
                Estimated ~{estimatedTransition}s per slide
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
                                {formatDuration(cfg.totalPlayTime)} total · {cfg.transitionTime}s per slide
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
