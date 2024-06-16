<script lang="ts">
    import { twMerge } from 'tailwind-merge';
    import type { Loading } from './types';

    import IconFileDamageFill from '~icons/ri/file-damage-fill';

    export let src: string = "";
    export let alt: string = "";
    export { className as class };
    export let loading: Loading = 'eager';
    export let height: number = 0;
    export let width: number = 0;
    export let hover: boolean = false;
    export let highlight: boolean = false;
    export let rounded: boolean = false;

    const videoExtensions = ['mp4', 'webm', 'ogg'];

    let className: string = '';
    let hasError = false;

    function isVideo(src: string): boolean {
        const extension = src.split('.').pop()?.toLowerCase();
        return videoExtensions.includes(extension ?? '');
    }

    function handleMediaError() {
        hasError = true;
    }

    $: buttonClasses = twMerge(
        'border-none bg-none p-0 cursor-pointer block placeholder overflow-hidden',
        hover ? 'opacity-90 transition-all duration-200 hover:opacity-100 hover:scale-105' : '',
        highlight ? 'ring-2 ring-primary-500 bg-initial' : '',
        rounded ? 'rounded-container-token' : '',
        height ? `h-${height}` : 'h-full',
        width ? `w-${width}` : 'w-full',
        className
    );

    $: mediaClasses = 'w-full';
    $: placeholderHeight = height ? `h-${height}` : 'h-32';
</script>

<button
    type="button"
    on:click
    on:dblclick
    on:keydown
    class={buttonClasses}
>
    {#if !src || hasError}
        <div class="flex items-center justify-center animate-pulse {placeholderHeight}">
            <div class="font-bold text-surface-500-400-token">
                {#if hasError}
                    <div><slot name="error"><IconFileDamageFill/></slot></div>
                {:else}
                    <div><slot name="not-found">???</slot></div>
                {/if}
            </div>
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
