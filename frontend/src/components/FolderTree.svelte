<script>
    import { createEventDispatcher } from 'svelte';
    import { collections, buildTree } from '../stores/collections.js';

    export let selectedIds = [];

    const dispatch = createEventDispatcher();

    function toggle(id) {
        if (selectedIds.includes(id)) {
            dispatch('change', selectedIds.filter(x => x !== id));
        } else {
            dispatch('change', [...selectedIds, id]);
        }
    }

    $: tree = buildTree($collections);
</script>

<div class="folder-tree">
    {#each tree as node}
        <label class="tree-item">
            <input
                type="checkbox"
                checked={selectedIds.includes(node.id)}
                on:change={() => toggle(node.id)}
            />
            <span>📂 {node.name}</span>
        </label>
        {#if node.children?.length > 0}
            <div class="tree-children">
                {#each node.children as child}
                    <label class="tree-item">
                        <input
                            type="checkbox"
                            checked={selectedIds.includes(child.id)}
                            on:change={() => toggle(child.id)}
                        />
                        <span>📄 {child.name}</span>
                    </label>
                {/each}
            </div>
        {/if}
    {/each}
</div>

<style>
    .folder-tree {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .tree-item {
        display: flex;
        align-items: center;
        gap: var(--space-sm);
        padding: 6px var(--space-sm);
        border-radius: var(--radius-sm);
        cursor: pointer;
        font-size: 0.875rem;
        color: var(--text-secondary);
        transition: background var(--duration-fast) var(--ease-out);
    }
    .tree-item:hover {
        background: var(--bg-hover);
        color: var(--text-primary);
    }

    .tree-item input[type="checkbox"] {
        accent-color: var(--accent);
    }

    .tree-children {
        padding-left: var(--space-lg);
    }
</style>
