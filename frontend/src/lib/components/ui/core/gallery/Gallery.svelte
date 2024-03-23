<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    import Media from '$lib/components/ui/core/media/Media.svelte';
    import type { MediaInfo } from '$lib/types';
    import { twMerge } from 'tailwind-merge';

    export let items: MediaInfo[] = [];
    export let group: string[] = [];
    export let multiple: boolean = false;

    let clickTimeout: NodeJS.Timeout;

    $: divClass = twMerge('grid', $$props.class);

    const dispatch = createEventDispatcher();

    function init(node: HTMLElement) {
        if (getComputedStyle(node).gap === 'normal') node.style.gap = 'inherit';
    }

    function handleClick(itemId: string) {
        toggleSelection(itemId);
    }

    function handleDoubleClick(event: MouseEvent, itemId: string) {
         // Ensure the item is selected on double-click
        if (!group.includes(itemId)) {
            group = multiple ? [...group, itemId] : [itemId];
        }

        dispatch('itemDoubleClick', { event: event, item: itemId});
    }

    function handleKeyDown(event: KeyboardEvent, itemId: string) {
        if (event.key === 'Enter' || event.key === ' ') {
            toggleSelection(itemId);
        }
    }

    function toggleSelection(itemId: string) {
        if (multiple) {
            const index = group.indexOf(itemId);
            if (index === -1) {
                group = [...group, itemId];
            } else {
                group = group.filter(id => id !== itemId);
            }
        } else {
            group = group[0] === itemId ? [] : [itemId];
        }
    }
</script>

<div {...$$restProps} class={divClass} use:init>
    {#each items as item}
        <slot {item}>
            <Media 
                on:click={() => handleClick(item.id)}
                on:dblclick={(e) => handleDoubleClick(e, item.id)}
                on:keydown={(e) => handleKeyDown(e, item.id)}
                src={item.src}
                alt={item.alt}
                title={item.title}
                selected={group.includes(item.id)}
                width="full"
                height="80" 
                interactable />
        </slot>
    {:else}
        <slot item={items[0]} />
    {/each}
</div>