<script>
    import { createEventDispatcher, onDestroy, onMount } from 'svelte';
    import { formatDuration, localFileUrl, localStreamUrl, needsTranscode } from '../services/utils.js';
    import { isPlaying, currentIndex, playerMedia } from '../stores/player.js';
    import Hls from 'hls.js';

    export let mediaList = [];
    export let timePerPicture = 5;
    export let soundSource = '';

    const dispatch = createEventDispatcher();

    let timer = null;
    let progress = 0;
    let progressTimer = null;
    let soundEl = null;

    $: currentMedia = mediaList[$currentIndex] || null;
    $: totalItems = mediaList.length;

    function startPlayback() {
        isPlaying.set(true);
        scheduleNext();
        startSound();
    }

    function pausePlayback() {
        isPlaying.set(false);
        clearTimers();
        pauseSound();
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

        const duration = timePerPicture * 1000;
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
            stopSound();
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

    // ─── Sound source management ───────────────────────────────
    function startSound() {
        if (!soundEl || !soundSource) return;
        soundEl.currentTime = 0;
        soundEl.play().catch(() => {}); // ignore autoplay rejection
    }

    function pauseSound() {
        if (!soundEl) return;
        soundEl.pause();
    }

    function stopSound() {
        if (!soundEl) return;
        soundEl.pause();
        soundEl.currentTime = 0;
    }

    function handleSoundEnded() {
        // If slideshow is still playing, loop the sound
        if ($isPlaying) {
            soundEl.currentTime = 0;
            soundEl.play().catch(() => {});
        }
    }

    // ─── Keyboard controls ─────────────────────────────────────
    function handleKeydown(e) {
        if (e.key === 'Escape') {
            dispatch('exit');
        } else if (e.key === 'ArrowRight' || e.key === ' ') {
            e.preventDefault();
            goNext();
        } else if (e.key === 'ArrowLeft') {
            e.preventDefault();
            goPrev();
        }
    }

    // Auto-start when component mounts
    onMount(() => {
        currentIndex.set(0);
        startPlayback();
    });

    onDestroy(() => {
        clearTimers();
        isPlaying.set(false);
        stopSound();
        destroyHls();
    });

    let hlsInstance = null;

    function attachHls(node) {
        destroyHls();
        if (currentMedia && currentMedia.type === 'video' && needsTranscode(currentMedia.path)) {
            const streamUrl = localStreamUrl(currentMedia.path);
            if (Hls.isSupported()) {
                hlsInstance = new Hls();
                hlsInstance.loadSource(streamUrl);
                hlsInstance.attachMedia(node);
                hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
                    node.play().catch(() => {});
                });
            }
        }
    }

    function destroyHls() {
        if (hlsInstance) {
            hlsInstance.destroy();
            hlsInstance = null;
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Hidden sound source audio element -->
{#if soundSource}
    <audio
        bind:this={soundEl}
        src={localFileUrl(soundSource)}
        on:ended={handleSoundEnded}
        preload="auto"
    ></audio>
{/if}

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
                {#if needsTranscode(currentMedia.path)}
                    <video
                        use:attachHls
                        on:ended={handleVideoEnded}
                        class="crossfade-enter"
                        controls
                    >
                        <track kind="captions" />
                    </video>
                {:else}
                    <video
                        src={localFileUrl(currentMedia.path)}
                        autoplay
                        on:ended={handleVideoEnded}
                        class="crossfade-enter"
                    >
                        <track kind="captions" />
                    </video>
                {/if}
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
            {timePerPicture}s/slide
        </span>

        {#if soundSource}
            <span class="player-sound-badge text-sm">🎵</span>
        {/if}

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

    .player-sound-badge {
        color: var(--accent);
        font-weight: 600;
    }
</style>
