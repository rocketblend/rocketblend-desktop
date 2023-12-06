<script lang="ts">
    import { twMerge } from 'tailwind-merge';
    import { onMount, onDestroy } from 'svelte';

    export let src: string = "";
    export let alt: string = "";
    export let title: string = "";

    export let selected: boolean = false;

    export let containerClass: string = "rounded-container-token";
    export let placeholderClass: string = "placeholder";
    export let loadingClass: string = "animate-pulse";

    export let height: string = "32";
    export let width: string = "32";

    export let OnClick: (event: MouseEvent) => void = () => {};
    export let OnKeyDown: (event: KeyboardEvent) => void = () => {};

    let mediaLoaded = false;

    const isWebm = (path: string) => path.endsWith(".webm");

    $: if (src || src == "") {mediaLoaded = false;}

    $: heightClass = `h-${height}`;
    $: widthClass = `w-${width}`;

    let mediaElement: HTMLElement | null = null;
    let mediaHeight: string = '';

    $: mediaClass = twMerge(containerClass, heightClass, widthClass, mediaLoaded ? 'h-auto' : '');
    $: holderClass = twMerge(mediaClass, placeholderClass, src !== "" && !mediaLoaded ? loadingClass : '');

    $: focusClass = selected ? "ring-2 ring-surface-50" : "";
    $: containerStyle = mediaHeight !== '' ? `height: ${mediaHeight};` : '';

    function onMediaLoad(event: Event) {
        mediaElement = event.target as HTMLElement;

        // setTimeout(() => {
        //     mediaLoaded = true;
        //     updateHeight();
        // }, Math.floor(Math.random() * 3000));

        mediaLoaded = true;
        updateHeight();
    }

    function updateHeight() {
        if (mediaElement) {
            if (mediaElement.offsetHeight != 0) {
                mediaHeight = `${mediaElement.offsetHeight}px`;
                return;
            }

            // If the height is 0, the image is not loaded yet, so we wait a bit and try again.
            setTimeout(updateHeight, 10);
        }
    }

    function onResize() {
        updateHeight();
    }

    onMount(() => {
        window.addEventListener('resize', onResize);
    });

    onDestroy(() => {
        window.removeEventListener('resize', onResize);
    });
</script>

<style>
    /* Hide overlay by default */
    .overlay {
        display: none;
    }

    /* Show overlay on hover */
    .hover-container:hover .overlay {
        display: flex;
    }
</style>

<div 
    class="hover-container {focusClass} {holderClass} relative inline-block overflow-hidden" 
    on:click={OnClick}
    on:keydown={OnKeyDown}
    role="button" 
    tabindex="0"
    style={containerStyle}
>
    {#if src}
        {#if isWebm(src)}
            <video
                class={mediaClass}
                class:hidden={!mediaLoaded}
                src={src} 
                autoplay loop muted playsinline 
                on:loadeddata={onMediaLoad}
            ></video>
        {:else}
            <img 
                class={mediaClass}
                class:hidden={!mediaLoaded}
                src={src} 
                alt={alt} 
                on:load={onMediaLoad}
            />
        {/if}
        <div class="overlay absolute inset-0 flex justify-center items-center text-white bg-primary-hover-token" class:hidden={!mediaLoaded}>
            <h6 class="font-bold">{title}</h6>
        </div>
    {:else}
        <div class="overlay flex justify-center items-center h-full bg-primary-hover-token">
            <h6 class="font-bold">{title}</h6>
        </div>
    {/if}
</div>
