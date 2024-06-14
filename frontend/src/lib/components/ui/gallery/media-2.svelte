<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { twMerge } from 'tailwind-merge';
    import type { Loading, MediaDetails } from './types';

    export let src: string;
    export let alt: string;
    export let className: string = '';
    export let loading: Loading = 'eager';
    export let hover: boolean = false;

    const dispatch = createEventDispatcher<{ click: MediaDetails }>();

    function handleClick() {
        dispatch('click', { src, alt, class: className, loading });
    }

    function isVideo(src: string): boolean {
        const videoExtensions = ['mp4', 'webm', 'ogg'];
        const extension = src.split('.').pop()?.toLowerCase();
        return videoExtensions.includes(extension ?? '');
    }

    $: buttonClasses = twMerge(
        'border-none bg-none p-0 cursor-pointer w-full block',
        className,
        hover ? 'opacity-90 transition-all duration-200 hover:opacity-100 hover:scale-105' : ''
    );

    $: mediaClasses = 'w-full';
</script>

<button
    type="button"
    on:click={handleClick}
    on:keydown={(e) => e.key === 'Enter' && handleClick()}
    class={buttonClasses}
>
    {#if isVideo(src)}
        <video src={src} class={mediaClasses} preload={loading === 'eager' ? 'auto' : 'metadata'} autoplay loop muted playsinline >
            Your browser does not support the video tag.
        </video>
    {:else}
        <img src={src} alt={alt} loading={loading} class={mediaClasses} />
    {/if}
</button>
