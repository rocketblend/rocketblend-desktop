<script lang="ts">
    import { twMerge } from 'tailwind-merge';

    export let src: string = "";
    export let alt: string = "";
    export let title: string = "";

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

    $: mediaClass = twMerge(containerClass, heightClass, widthClass, mediaLoaded ? 'h-auto' : '');
    $: holderClass = twMerge(mediaClass, placeholderClass, src !== "" && !mediaLoaded ? loadingClass : '');

    function onMediaLoad() {
        // setTimeout(() => {
        //     mediaLoaded = true;
        // }, Math.floor(Math.random() * 3000));
        
        mediaLoaded = true;
    }
</script>


<div 
    class="focus:ring-2 focus:ring-surface-50 {holderClass}" 
    on:click={OnClick}
    on:keydown={OnKeyDown}
    role="button" 
    tabindex="0"
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
    {:else}
        <div class="{holderClass}"></div>
    {/if}
</div>