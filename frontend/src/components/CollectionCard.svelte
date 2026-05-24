<script>
    import { onMount } from 'svelte';
    import * as api from '../services/api.js';
    import { localThumbUrl } from '../services/utils.js';
    import { lazyload } from '../actions/lazyload.js';

    export let collection = {};
    export let onClick = () => {};

    let thumbnail = '';
    let mediaCount = 0;
    let subCount = 0;

    $: thumbSrc = thumbnail ? localThumbUrl(thumbnail) : '';

    onMount(async () => {
        try {
            const [thumb, count, subs] = await Promise.all([
                api.getCollectionThumbnail(collection.id),
                api.getCollectionMediaCount(collection.id),
                api.getSubCollections(collection.id),
            ]);
            thumbnail = thumb || '';
            mediaCount = count || 0;
            subCount = subs ? subs.length : 0;
        } catch (e) {
            console.error('CollectionCard load error:', e);
        }
    });
</script>

<button class="collection-card" on:click={() => onClick(collection)}>
    <div class="card-cover">
        {#if thumbSrc}
            <img use:lazyload={thumbSrc} alt={collection.name} class="cover-image" />
        {:else}
            <div class="cover-placeholder">
                <span class="cover-icon">📁</span>
            </div>
        {/if}
        <div class="card-overlay">
            {#if subCount > 0}
                <span class="badge badge-sub">📂 {subCount}</span>
            {/if}
        </div>
    </div>
    <div class="card-body">
        <div class="card-title truncate" title={collection.name}>{collection.name}</div>
        <div class="card-subtitle text-sm text-muted">
            {mediaCount} file{mediaCount !== 1 ? 's' : ''}
        </div>
    </div>
</button>

<style>
    .collection-card {
        background: var(--bg-tertiary);
        border: 1px solid var(--glass-border);
        border-radius: var(--radius-lg);
        overflow: hidden;
        cursor: pointer;
        transition: all var(--duration-normal) var(--ease-out);
        text-align: left;
        font-family: var(--font-sans);
        color: var(--text-primary);
        padding: 0;
        width: 100%;
    }
    .collection-card:hover {
        border-color: rgba(255,255,255,0.12);
        transform: translateY(-3px);
        box-shadow: var(--shadow-lg);
    }
    .collection-card:hover .cover-image {
        transform: scale(1.06);
    }

    .card-cover {
        position: relative;
        aspect-ratio: 16/10;
        overflow: hidden;
        background: var(--bg-primary);
    }

    .cover-image {
        width: 100%;
        height: 100%;
        object-fit: cover;
        object-position: center;
        display: block;
        transition: transform var(--duration-slow) var(--ease-out);
    }

    .cover-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, var(--bg-tertiary) 0%, var(--bg-elevated) 100%);
    }
    .cover-icon {
        font-size: 3rem;
        opacity: 0.4;
    }

    .card-overlay {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        padding: 6px 8px;
        display: flex;
        gap: 4px;
        justify-content: flex-end;
    }

    .badge-sub {
        background: rgba(0, 0, 0, 0.6);
        color: #fff;
        font-size: 0.65rem;
        padding: 2px 6px;
        border-radius: var(--radius-full);
        backdrop-filter: blur(4px);
    }

    .card-body {
        padding: var(--space-sm) var(--space-md);
    }

    .card-title {
        font-size: 0.9rem;
        font-weight: 600;
        line-height: 1.3;
    }

    .card-subtitle {
        margin-top: 2px;
    }
</style>
