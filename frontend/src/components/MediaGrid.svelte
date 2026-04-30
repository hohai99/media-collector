<script>
    import MediaCard from './MediaCard.svelte';

    export let items = [];
</script>

{#if items.length === 0}
    <div class="empty-state">
        <span class="icon">📭</span>
        <h3>No media files</h3>
        <p class="text-sm text-muted">Scan your master folder to discover media files</p>
    </div>
{:else}
    <div class="media-grid">
        {#each items as item (item.id)}
            <MediaCard media={item} />
        {/each}
    </div>
{/if}
