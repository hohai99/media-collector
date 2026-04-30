<script>
    import { createEventDispatcher, onDestroy } from 'svelte';
    import { formatDuration, localFileUrl } from '../services/utils.js';
    import { isPlaying, currentIndex, playerMedia } from '../stores/player.js';

    export let mediaList = [];
    export let transitionTime = 5;

    const dispatch = createEventDispatcher();

    let timer = null;
    let progress = 0;
    let progressTimer = null;

    $: currentMedia = mediaList[$currentIndex] || null;
    $: totalItems = mediaList.length;

    function startPlayback() {
        isPlaying.set(true);
        scheduleNext();
    }

    function pausePlayback() {
        isPlaying.set(false);
        clearTimers();
    }

    function togglePlayback() {
        if ($isPlaying) pausePlayback();
        else startPlayback();
    }

    function scheduleNext() {
        clearTimers();
        progress = 0;

        if (currentMedia?.type === 'video' || currentMedia?.type === 'audio') {
            // Video/audio plays for full duration — no timer needed
            return;
        }

        const duration = transitionTime * 1000;
        const interval = 50;
        progressTimer = setInterval(() => {
            progress += (interval / duration) * 100;
            if (progress >= 100) progress = 100;
        }, interval);

        timer = setTimeout(() => {
            goNext();
        }, duration);
    }

    function goNext() {
        if ($currentIndex < totalItems - 1) {
            currentIndex.update(i => i + 1);
            if ($isPlaying) scheduleNext();
        } else {
            // End of playlist
            isPlaying.set(false);
            clearTimers();
        }
    }

    function goPrev() {
        if ($currentIndex > 0) {
            currentIndex.update(i => i - 1);
            if ($isPlaying) scheduleNext();
        }
    }

    function handleVideoEnded() {
        goNext();
    }

    function clearTimers() {
        if (timer) { clearTimeout(timer); timer = null; }
        if (progressTimer) { clearInterval(progressTimer); progressTimer = null; }
        progress = 0;
    }

    function handleKeydown(e) {
        if (e.key === 'Escape') {
            dispatch('exit');
        } else if (e.key === 'F11') {
            e.preventDefault();
            // Toggle fullscreen is handled by the player page
        } else if (e.key === 'ArrowRight' || e.key === ' ') {
            e.preventDefault();
            goNext();
        } else if (e.key === 'ArrowLeft') {
            e.preventDefault();
            goPrev();
        }
    }

    // Auto-start when component mounts
    import { onMount } from 'svelte';
    onMount(() => {
        currentIndex.set(0);
        startPlayback();
    });

    onDestroy(() => {
        clearTimers();
        isPlaying.set(false);
    });
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="player-fullscreen">
    {#if currentMedia}
        {#if currentMedia.type === 'image'}
            {#key $currentIndex}
                <img
                    class="crossfade-enter"
                    src={localFileUrl(currentMedia.path)}
                    alt={currentMedia.name}
                />
            {/key}
        {:else if currentMedia.type === 'video'}
            {#key $currentIndex}
                <video
                    src={localFileUrl(currentMedia.path)}
                    autoplay
                    on:ended={handleVideoEnded}
                    class="crossfade-enter"
                >
                    <track kind="captions" />
                </video>
            {/key}
        {:else if currentMedia.type === 'audio'}
            <div class="audio-display">
                <span class="audio-icon">🎵</span>
                <h2>{currentMedia.name}</h2>
                {#key $currentIndex}
                    <audio
                        src={localFileUrl(currentMedia.path)}
                        autoplay
                        on:ended={handleVideoEnded}
                    />
                {/key}
            </div>
        {/if}
    {/if}

    <div class="player-controls">
        <button class="btn btn-ghost" on:click={goPrev} disabled={$currentIndex === 0}>⏮</button>
        <button class="btn btn-ghost" on:click={togglePlayback}>
            {$isPlaying ? '⏸' : '▶️'}
        </button>
        <button class="btn btn-ghost" on:click={goNext} disabled={$currentIndex >= totalItems - 1}>⏭</button>

        <div class="player-progress">
            <div class="player-progress-fill" style="width: {progress}%"></div>
        </div>

        <span class="player-counter text-sm">
            {$currentIndex + 1} / {totalItems}
        </span>

        <span class="player-time text-sm text-muted">
            {formatDuration(transitionTime)}/slide
        </span>

        <button class="btn btn-ghost btn-sm" on:click={() => dispatch('exit')}>✕ Exit</button>
    </div>
</div>

<style>
    .audio-display {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--space-lg);
        color: white;
    }
    .audio-icon {
        font-size: 6rem;
    }
    .audio-display h2 {
        color: var(--text-secondary);
    }

    .player-counter {
        color: white;
        font-weight: 500;
    }

    .player-controls button:disabled {
        opacity: 0.3;
        cursor: not-allowed;
    }
</style>
