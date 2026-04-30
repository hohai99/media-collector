<script>
    import { onMount } from 'svelte';
    import Sidebar from './components/Sidebar.svelte';
    import Toast from './components/Toast.svelte';
    import Home from './pages/Home.svelte';
    import Import from './pages/Import.svelte';
    import PlayerBuilder from './pages/PlayerBuilder.svelte';
    import Player from './pages/Player.svelte';
    import { currentPage, masterFolder, showToast } from './stores/app.js';
    import { collections } from './stores/collections.js';
    import * as api from './services/api.js';

    onMount(async () => {
        try {
            const folder = await api.getMasterFolder();
            if (folder) {
                masterFolder.set(folder);
                // Sync filesystem → DB and scan for new media
                await api.syncAndScan();
                const tree = await api.getCollectionTree();
                collections.set(tree || []);
            }
        } catch (e) {
            console.error('Init error:', e);
        }
    });
</script>

<Toast />

{#if $currentPage === 'player'}
    <Player />
{:else}
    <div class="app-layout">
        <Sidebar />
        <main class="main-content">
            {#if $currentPage === 'home'}
                <Home />
            {:else if $currentPage === 'import'}
                <Import />
            {:else if $currentPage === 'player-builder'}
                <PlayerBuilder />
            {/if}
        </main>
    </div>
{/if}
