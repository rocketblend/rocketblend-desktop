<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { twMerge } from 'tailwind-merge';
    import type { Loading, ImageDetails } from './types';

    export let src: string;
    export let alt: string;
    export let className: string = '';
    export let loading: Loading = 'eager';
    export let hover: boolean = false;

    const dispatch = createEventDispatcher<{ click: ImageDetails }>();

    function handleClick() {
        dispatch('click', { src, alt, class: className, loading });
    }

    $: buttonClasses = twMerge(
        'border-none bg-none p-0 cursor-pointer w-full block',
        className,
        hover ? 'opacity-90 transition-all duration-200 hover:opacity-100 hover:scale-105' : ''
    );

    $: imgClasses = 'w-full';
</script>

<button
    type="button"
    on:click={handleClick}
    on:keydown={(e) => e.key === 'Enter' && handleClick()}
    class={buttonClasses}
>
    <img src={src} alt={alt} loading={loading} class={imgClasses} />
</button>
