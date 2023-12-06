<script lang="ts">
	import Media from '$lib/components/core/media/Media.svelte';
    import type { ImgType } from '$lib/types';
    import { twMerge } from 'tailwind-merge';

    export let items: ImgType[] = [];

    $: divClass = twMerge('grid', $$props.class);

    function init(node: HTMLElement) {
        if (getComputedStyle(node).gap === 'normal') node.style.gap = 'inherit';
    }
    
</script>
  
<div {...$$restProps} class={divClass} use:init>
    {#each items as item}
        <slot {item}>
        <Media src={item.src.img.src} alt={item.alt} width="full" height="64"/>
        </slot>
    {:else}
        <slot item={items[0]} />
    {/each}
</div>