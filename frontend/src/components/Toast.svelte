<script>
    import { toasts } from '../stores/app.js';
</script>

{#if $toasts.length > 0}
    <div class="toast-container">
        {#each $toasts as toast (toast.id)}
            <div class="toast {toast.type}">
                {#if toast.type === 'success'}✅
                {:else if toast.type === 'error'}❌
                {:else}ℹ️
                {/if}
                <span>{toast.message}</span>
            </div>
        {/each}
    </div>
{/if}
