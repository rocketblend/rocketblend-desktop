<script lang="ts">
    import Media from '$lib/components/core/media/Media.svelte';
    import type { MediaInfo } from '$lib/types';
    import { twMerge } from 'tailwind-merge';

    export let items: MediaInfo[] = [];
    export let group: string[] = [];
    export let multiple: boolean = false;

    $: divClass = twMerge('grid', $$props.class);

    function init(node: HTMLElement) {
        if (getComputedStyle(node).gap === 'normal') node.style.gap = 'inherit';
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

    function handleKeyDown(event: KeyboardEvent, itemId: string) {
        if (event.key === 'Enter' || event.key === ' ') {
            toggleSelection(itemId);
        }
    }
</script>

<div {...$$restProps} class={divClass} use:init>
    {#each items as item}
        <slot {item}>
            <Media 
                OnClick={() => toggleSelection(item.id)}
                OnKeyDown={(e) => handleKeyDown(e, item.id)}
                src={item.src}
                alt={item.alt}
                title={item.title}
                width="full"
                height="80" />
        </slot>
    {:else}
        <slot item={items[0]} />
    {/each}
</div>