<script lang="ts">
    import type { ImgType } from '$lib/types';
    import { twMerge } from 'tailwind-merge';

    import Img from '@zerodevx/svelte-img';

    export let items: ImgType[] = [];
    export let imgClass: string = 'h-auto max-w-full rounded-container-token image-fade-in';

    $: divClass = twMerge('grid', $$props.class);

    function init(node: HTMLElement) {
        if (getComputedStyle(node).gap === 'normal') node.style.gap = 'inherit';
    }
    
</script>
  
<div {...$$restProps} class={divClass} use:init>
    {#each items as item}
        <slot {item}>
        <div>
            <img class={twMerge(imgClass, $$props.classImg)} src={item.src.img.src} alt={item.alt} />
            <!-- <Img class={twMerge(imgClass, $$props.classImg)} src={item.src} alt={item.alt} /> -->
        </div>
        </slot>
    {:else}
        <slot item={items[0]} />
    {/each}
</div>