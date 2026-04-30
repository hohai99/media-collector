<script>
    import { currentPage, masterFolder, navigateTo, showToast } from '../stores/app.js';
    import { collections, selectedCollection } from '../stores/collections.js';
    import { buildTree } from '../stores/collections.js';
    import * as api from '../services/api.js';

    let showNewCollection = false;
    let newCollectionName = '';

    const navItems = [
        { id: 'home',           icon: '📁', label: 'Collections' },
        { id: 'import',         icon: '📥', label: 'Import' },
        { id: 'player-builder', icon: '🎛️', label: 'Player' },
    ];

    async function handleCreateCollection() {
        if (!newCollectionName.trim()) return;
        try {
            await api.createCollection(newCollectionName.trim());
            showToast('Collection created', 'success');
            newCollectionName = '';
            showNewCollection = false;
            await refreshCollections();
        } catch (e) {
            showToast(e, 'error');
        }
    }

    async function refreshCollections() {
        try {
            const tree = await api.getCollectionTree();
            collections.set(tree || []);
        } catch (e) {
            console.error('Failed to load collections:', e);
        }
    }

    function selectCol(col) {
        selectedCollection.set(col);
        navigateTo('home');
    }

    $: collTree = buildTree($collections);
</script>

<aside class="sidebar">
    <div class="sidebar-brand">
        <span class="brand-icon">💎</span>
        <span class="brand-text">Media Collector</span>
    </div>

    <nav class="sidebar-nav">
        {#each navItems as item}
            <button
                class="nav-item"
                class:active={$currentPage === item.id}
                on:click={() => navigateTo(item.id)}
            >
                <span class="nav-icon">{item.icon}</span>
                <span class="nav-label">{item.label}</span>
            </button>
        {/each}
    </nav>

    <div class="sidebar-section">
        <div class="section-header">
            <span class="section-title">Collections</span>
            <button class="btn btn-ghost btn-icon btn-sm" on:click={() => showNewCollection = !showNewCollection} title="New collection">
                ＋
            </button>
        </div>

        {#if showNewCollection}
            <div class="new-collection-form">
                <input
                    class="input"
                    bind:value={newCollectionName}
                    placeholder="Collection name…"
                    on:keydown={(e) => e.key === 'Enter' && handleCreateCollection()}
                />
            </div>
        {/if}

        <div class="collection-list">
            {#if collTree.length === 0}
                <div class="text-muted text-sm" style="padding: 8px 12px;">No collections yet</div>
            {:else}
                {#each collTree as node}
                    <button
                        class="collection-item"
                        class:active={$selectedCollection?.id === node.id}
                        on:click={() => selectCol(node)}
                    >
                        <span>📂</span>
                        <span class="truncate">{node.name}</span>
                    </button>
                    {#if node.children?.length > 0}
                        <div class="collection-children">
                            {#each node.children as child}
                                <button
                                    class="collection-item child"
                                    class:active={$selectedCollection?.id === child.id}
                                    on:click={() => selectCol(child)}
                                >
                                    <span>📄</span>
                                    <span class="truncate">{child.name}</span>
                                </button>
                            {/each}
                        </div>
                    {/if}
                {/each}
            {/if}
        </div>
    </div>

    {#if $masterFolder}
        <div class="sidebar-footer">
            <div class="text-sm text-muted truncate" title={$masterFolder}>
                📍 {$masterFolder}
            </div>
        </div>
    {/if}
</aside>

<style>
    .sidebar {
        width: var(--sidebar-width);
        min-width: var(--sidebar-width);
        height: 100vh;
        background: var(--bg-secondary);
        border-right: 1px solid var(--glass-border);
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .sidebar-brand {
        padding: var(--space-lg) var(--space-md);
        display: flex;
        align-items: center;
        gap: var(--space-sm);
        border-bottom: 1px solid var(--glass-border);
    }
    .brand-icon { font-size: 1.5rem; }
    .brand-text {
        font-weight: 700;
        font-size: 1rem;
        background: linear-gradient(135deg, var(--accent-start), var(--accent-end));
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
    }

    .sidebar-nav {
        padding: var(--space-sm);
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .nav-item {
        display: flex;
        align-items: center;
        gap: var(--space-sm);
        padding: var(--space-sm) var(--space-md);
        border: none;
        background: transparent;
        color: var(--text-secondary);
        border-radius: var(--radius-md);
        cursor: pointer;
        font-family: var(--font-sans);
        font-size: 0.875rem;
        transition: all var(--duration-fast) var(--ease-out);
        width: 100%;
        text-align: left;
    }
    .nav-item:hover {
        background: var(--bg-hover);
        color: var(--text-primary);
    }
    .nav-item.active {
        background: var(--accent-muted);
        color: var(--text-accent);
    }

    .sidebar-section {
        flex: 1;
        overflow-y: auto;
        padding: var(--space-sm);
        border-top: 1px solid var(--glass-border);
    }

    .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: var(--space-xs) var(--space-sm);
        margin-bottom: var(--space-xs);
    }
    .section-title {
        font-size: 0.7rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--text-muted);
    }

    .new-collection-form {
        padding: 0 var(--space-sm) var(--space-sm);
    }

    .collection-list {
        display: flex;
        flex-direction: column;
        gap: 1px;
    }

    .collection-item {
        display: flex;
        align-items: center;
        gap: var(--space-sm);
        padding: 6px var(--space-sm);
        border: none;
        background: transparent;
        color: var(--text-secondary);
        border-radius: var(--radius-sm);
        cursor: pointer;
        font-family: var(--font-sans);
        font-size: 0.825rem;
        transition: all var(--duration-fast) var(--ease-out);
        width: 100%;
        text-align: left;
    }
    .collection-item:hover {
        background: var(--bg-hover);
        color: var(--text-primary);
    }
    .collection-item.active {
        background: var(--accent-muted);
        color: var(--text-accent);
    }

    .collection-children {
        padding-left: var(--space-lg);
    }
    .collection-item.child {
        font-size: 0.8rem;
    }

    .sidebar-footer {
        padding: var(--space-sm) var(--space-md);
        border-top: 1px solid var(--glass-border);
    }
</style>
