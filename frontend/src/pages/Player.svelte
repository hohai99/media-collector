<script>
    import PlayerView from '../components/PlayerView.svelte';
    import { activeConfig, playerMedia } from '../stores/player.js';
    import { navigateTo } from '../stores/app.js';

    function handleExit() {
        navigateTo('player-builder');
    }

    $: timePerPicture = $activeConfig?.timePerPicture || $activeConfig?.transitionTime || 5;
    $: soundSource = $activeConfig?.soundSource || '';
</script>

{#if $playerMedia.length > 0}
    <PlayerView
        mediaList={$playerMedia}
        {timePerPicture}
        {soundSource}
        on:exit={handleExit}
    />
{:else}
    <div class="page-container">
        <div class="empty-state">
            <span class="icon">🎬</span>
            <h3>No Media to Play</h3>
            <p class="text-sm text-muted">Go back and select a player config to start playback.</p>
            <button class="btn btn-primary" on:click={() => navigateTo('player-builder')}>
                ← Back to Player Builder
            </button>
        </div>
    </div>
{/if}
