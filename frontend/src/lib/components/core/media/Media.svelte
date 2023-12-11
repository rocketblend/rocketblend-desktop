<script lang="ts">
    import { twMerge } from 'tailwind-merge';
    import { isVideo } from '$lib/components/utils';

    export let src: string = "";
    export let alt: string = "";
    export let title: string = "";
    export let selected: boolean = false;
    export let containerClass: string = "rounded-container-token";
    export let placeholderClass: string = "placeholder";
    export let loadingClass: string = "animate-pulse";
    export let height: string = "32";
    export let width: string = "32";
    export let simulateLatency: boolean = false;

    let mediaLoaded = false;

    $: heightClass = `h-${height}`;
    $: widthClass = `w-${width}`;
    $: mediaClass = twMerge(containerClass, heightClass, widthClass, mediaLoaded ? 'h-fit' : '');
    $: holderClass = twMerge(mediaClass, placeholderClass, src && !mediaLoaded ? loadingClass : '');
    $: focusClass = selected ? "ring-2 ring-surface-50" : "";

    function onMediaLoad() {
        if (!simulateLatency) {
            mediaLoaded = true;
            return;
        }

        setTimeout(() => {
            mediaLoaded = true;
        }, Math.floor(Math.random() * 5000));
    }
</script>

<div 
    class="group {focusClass} {holderClass} relative inline-block overflow-hidden" 
    on:click
    on:dblclick
    on:keydown
    role="button" 
    tabindex="0"
>
    {#if src}
        {#if isVideo(src)}
            <video
                class="{mediaClass} {mediaLoaded ? 'block' : 'hidden'}"
                src={src} 
                autoplay loop muted playsinline 
                on:loadeddata={onMediaLoad}
            ></video>
        {:else}
            <img 
                class="{mediaClass} {mediaLoaded ? 'block' : 'hidden'}"
                src={src} 
                alt={alt} 
                on:load={onMediaLoad}
            />
        {/if}
        <div class="overlay absolute inset-0 flex justify-center items-center text-white bg-primary-hover-token group-hover:flex hidden">
            <h6 class="font-bold">{title}</h6>
        </div>
    {:else}
        <div class="overlay flex justify-center items-center h-full bg-primary-hover-token group-hover:flex hidden">
            <h6 class="font-bold">{title}</h6>
        </div>
    {/if}
</div>