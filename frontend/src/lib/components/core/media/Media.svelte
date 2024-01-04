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
    export let interactable: boolean = false;

    let mediaLoaded: boolean = false;

    $: if (src || src == "") {mediaLoaded = false;}
    $: heightClass = `h-${height}`;
    $: widthClass = `w-${width}`;
    $: mediaClass = twMerge(containerClass, heightClass, widthClass, mediaLoaded ? 'h-fit' : '');
    $: holderClass = twMerge(mediaClass, placeholderClass, src && !mediaLoaded ? loadingClass : '');
    $: focusClass = selected ? "ring-2 ring-primary-500 bg-initial" : "";
    $: cursorClass = interactable ? 'cursor-pointer' : 'cursor-default';

    function onMediaLoad() {
        if (simulateLatency) {
            setTimeout(() => { mediaLoaded = true; }, randomDelay());
        } else {
            mediaLoaded = true;
        }
    }

    function randomDelay() {
        return Math.floor(Math.random() * 5000);
    }
</script>

<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
<div 
    class={`group ${interactable ? focusClass : ''} ${holderClass} ${cursorClass} relative inline-block overflow-hidden`} 
    on:click
    on:dblclick
    on:keydown
    role={interactable ? 'button' : undefined}
    tabindex={interactable ? 0 : undefined}
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
            {#if interactable}
                <div class="overlay absolute inset-0 flex justify-center items-center text-white hover:variant-glass-surface group-hover:flex hidden">
                    <h6 class="font-bold">{title}</h6>
                </div>
            {/if}
    {:else}
        {#if interactable}
            <div class="overlay flex justify-center items-center h-full hover:variant-glass-surface group-hover:flex hidden">
                <h6 class="font-bold">{title}</h6>
            </div>
        {/if}
    {/if}
</div>