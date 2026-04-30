<script>
    import { selectedMedia, toggleMediaSelection } from '../stores/media.js';
    import { formatSize, mediaTypeIcon, localFileUrl } from '../services/utils.js';

    export let media = {};

    $: isSelected = $selectedMedia.includes(media.id);

    let showPreview = false;

    function handleClick(e) {
        // If shift/ctrl held, toggle selection; otherwise open preview
        if (e.ctrlKey || e.shiftKey || e.metaKey) {
            toggleMediaSelection(media.id);
        } else {
            showPreview = true;
        }
    }

    function handleSelect(e) {
        e.stopPropagation();
        toggleMediaSelection(media.id);
    }

    function closePreview() {
        showPreview = false;
    }
</script>

<div
    class="media-card selectable"
    class:selected={isSelected}
    on:click={handleClick}
    on:keydown={(e) => e.key === 'Enter' && handleClick(e)}
    role="button"
    tabindex="0"
>
    <div class="card-preview">
        {#if media.type === 'image'}
            <div class="preview-image" style="background-image: url('{localFileUrl(media.path)}')"></div>
        {:else}
            <div class="preview-placeholder">
                <span class="preview-icon">{mediaTypeIcon(media.type)}</span>
            </div>
        {/if}
        <span class="badge badge-{media.type}">{media.type}</span>
        <button class="select-btn" class:selected={isSelected} on:click={handleSelect} title="Select">
            {isSelected ? '✓' : ''}
        </button>
    </div>
    <div class="card-info">
        <div class="card-name truncate" title={media.name}>{media.name}</div>
        <div class="card-meta text-sm text-muted">{formatSize(media.size)}</div>
    </div>
</div>

<!-- Fullscreen Preview Modal -->
{#if showPreview}
    <div class="preview-overlay" on:click={closePreview} on:keydown={(e) => e.key === 'Escape' && closePreview()} role="dialog" tabindex="-1">
        <button class="preview-close" on:click={closePreview}>✕</button>
        <div class="preview-content" on:click|stopPropagation role="presentation">
            {#if media.type === 'image'}
                <img src={localFileUrl(media.path)} alt={media.name} />
            {:else if media.type === 'video'}
                <!-- svelte-ignore a11y-media-has-caption -->
                <video src={localFileUrl(media.path)} controls autoplay>
                    <track kind="captions" />
                </video>
            {:else if media.type === 'audio'}
                <div class="audio-preview">
                    <span class="audio-big-icon">🎵</span>
                    <h3>{media.name}</h3>
                    <audio src={localFileUrl(media.path)} controls autoplay></audio>
                </div>
            {/if}
        </div>
        <div class="preview-info">
            <span>{media.name}</span>
            <span class="text-muted">·</span>
            <span class="text-muted">{formatSize(media.size)}</span>
        </div>
    </div>
{/if}

<style>
    .media-card {
        background: var(--bg-tertiary);
        border: 1px solid var(--glass-border);
        border-radius: var(--radius-md);
        overflow: hidden;
        transition: all var(--duration-fast) var(--ease-out);
        cursor: pointer;
    }
    .media-card:hover {
        border-color: rgba(255,255,255,0.1);
        transform: translateY(-2px);
        box-shadow: var(--shadow-md);
    }
    .media-card.selected {
        border-color: var(--accent);
        box-shadow: 0 0 0 2px var(--accent-muted);
    }

    .card-preview {
        position: relative;
        aspect-ratio: 4/3;
        overflow: hidden;
        background: var(--bg-primary);
    }

    .preview-image {
        width: 100%;
        height: 100%;
        background-size: cover;
        background-position: center;
        transition: transform var(--duration-normal) var(--ease-out);
    }
    .media-card:hover .preview-image {
        transform: scale(1.05);
    }

    .preview-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, var(--bg-tertiary), var(--bg-primary));
    }
    .preview-icon {
        font-size: 2.5rem;
        opacity: 0.6;
    }

    .card-preview .badge {
        position: absolute;
        bottom: 6px;
        left: 6px;
    }

    .select-btn {
        position: absolute;
        top: 6px;
        right: 6px;
        width: 22px;
        height: 22px;
        border-radius: var(--radius-full);
        border: 2px solid rgba(255,255,255,0.4);
        background: rgba(0,0,0,0.3);
        color: #fff;
        font-size: 0.65rem;
        font-weight: 700;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        opacity: 0;
        transition: all var(--duration-fast) var(--ease-out);
        z-index: 2;
        backdrop-filter: blur(4px);
    }
    .media-card:hover .select-btn,
    .select-btn.selected {
        opacity: 1;
    }
    .select-btn.selected {
        background: var(--accent);
        border-color: var(--accent);
    }

    .card-info {
        padding: var(--space-sm);
    }

    .card-name {
        font-size: 0.825rem;
        font-weight: 500;
    }

    .card-meta {
        margin-top: 2px;
    }

    /* Preview overlay */
    .preview-overlay {
        position: fixed;
        inset: 0;
        z-index: 30000;
        background: rgba(0, 0, 0, 0.92);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        animation: fade-in var(--duration-fast) var(--ease-out);
    }

    .preview-close {
        position: absolute;
        top: 16px;
        right: 16px;
        width: 36px;
        height: 36px;
        border-radius: var(--radius-full);
        border: none;
        background: rgba(255,255,255,0.1);
        color: #fff;
        font-size: 1rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background var(--duration-fast) var(--ease-out);
        z-index: 2;
    }
    .preview-close:hover {
        background: rgba(255,255,255,0.2);
    }

    .preview-content {
        max-width: 90vw;
        max-height: 80vh;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .preview-content img {
        max-width: 90vw;
        max-height: 80vh;
        object-fit: contain;
        border-radius: var(--radius-md);
    }
    .preview-content video {
        max-width: 90vw;
        max-height: 80vh;
        object-fit: contain;
        border-radius: var(--radius-md);
    }

    .audio-preview {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--space-lg);
        color: white;
    }
    .audio-big-icon { font-size: 5rem; }
    .audio-preview h3 { color: var(--text-secondary); }

    .preview-info {
        margin-top: var(--space-md);
        display: flex;
        gap: var(--space-sm);
        font-size: 0.875rem;
        color: var(--text-primary);
    }
</style>
