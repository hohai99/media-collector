<script>
    export let title = '';
    export let show = false;

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    function close() {
        dispatch('close');
    }

    function handleKeydown(e) {
        if (e.key === 'Escape') close();
    }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
    <div class="modal-overlay" on:click|self={close}>
        <div class="modal-content">
            <div class="modal-header">
                <h3>{title}</h3>
                <button class="btn btn-ghost btn-icon" on:click={close}>✕</button>
            </div>
            <slot />
        </div>
    </div>
{/if}
