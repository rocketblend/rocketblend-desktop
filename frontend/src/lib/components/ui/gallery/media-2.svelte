<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { twMerge } from 'tailwind-merge';
    import type { Loading, MediaDetails } from './types';

    export let src: string = "";
    export let alt: string;
    export let className: string = '';
    export let loading: Loading = 'eager';
    export let hover: boolean = false;
    export let highlight: boolean = false;
    export let rounded: boolean = false;

    const dispatch = createEventDispatcher<{ click: MediaDetails }>();
    const videoExtensions = ['mp4', 'webm', 'ogg'];

    let hasError = false;

    function handleClick() {
        dispatch('click', { src, alt, class: className, highlight });
    }

    function isVideo(src: string): boolean {
        const extension = src.split('.').pop()?.toLowerCase();
        return videoExtensions.includes(extension ?? '');
    }

    function handleMediaError() {
        hasError = true;
    }

    $: buttonClasses = twMerge(
        'border-none bg-none p-0 cursor-pointer w-full h-full block placeholder overflow-hidden',
        className,
        hover ? 'opacity-90 transition-all duration-200 hover:opacity-100 hover:scale-105' : '',
        highlight ? 'ring-2 ring-primary-500 bg-initial' : '',
        rounded ? 'rounded-container-token' : ''
    );

    $: mediaClasses = 'w-full';
</script>

<button
    type="button"
    on:click={handleClick}
    on:keydown={(e) => e.key === 'Enter' && handleClick()}
    class={buttonClasses}
>
    {#if !src || hasError}
        <div class="flex items-center justify-center animate-pulse h-32">
            <span class="font-bold text-surface-500-400-token">
                {#if hasError}
                    Failed to load media
                {:else}
                    No media found
                {/if}
            </span>
        </div>
    {:else}
        {#if isVideo(src)}
            <video 
                src={src} 
                class={mediaClasses} 
                preload={loading === 'eager' ? 'auto' : 'metadata'} 
                autoplay 
                loop 
                muted 
                playsinline 
                on:error={handleMediaError}
                on:stalled={handleMediaError}
            >
                Your browser does not support the video tag.
            </video>
        {:else}
            <img 
                src={src} 
                alt={alt} 
                loading={loading} 
                class={mediaClasses} 
                on:error={handleMediaError} 
            />
        {/if}
    {/if}
</button>
