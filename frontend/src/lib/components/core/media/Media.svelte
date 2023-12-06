<script lang="ts">
    import { twMerge } from 'tailwind-merge';

    export let src: string = "";
    export let alt: string = "";

    export let containerClass: string = "rounded-container-token";
    export let placeholderClass: string = "placeholder";
    export let loadingClass: string = "animate-pulse";

    export let height: string = "32";
    export let width: string = "32";

    let mediaLoaded = false;

    const isWebm = (path: string) => path.endsWith(".webm");
    
    // const gradients = [
    //     "variant-gradient-primary-secondary",
    //     "variant-gradient-secondary-tertiary",
    //     "variant-gradient-tertiary-primary",
    //     "variant-gradient-secondary-primary",
    //     "variant-gradient-tertiary-secondary",
    //     "variant-gradient-primary-tertiary",
    //     "variant-gradient-success-warning",
    //     "variant-gradient-warning-error",
    //     "variant-gradient-error-success",
    //     "variant-gradient-warning-success",
    //     "variant-gradient-error-warning",
    //     "variant-gradient-success-error",
    // ];
    // const randomGradient = gradients[Math.floor(Math.random() * gradients.length)];
    // $: placeholderGradientClass = src === "" ? `bg-gradient-to-b ${randomGradient}` : "animate-pulse";

    $: if (src || src == "") {mediaLoaded = false;}

    $: heightClass = `h-${height}`;
    $: widthClass = `w-${width}`;

    $: mediaClass = twMerge(containerClass, heightClass, widthClass);
    $: holderClass = twMerge(mediaClass, placeholderClass, src !== "" && !mediaLoaded ? loadingClass : '');

    function onMediaLoad() {
        // setTimeout(() => {
        //     mediaLoaded = true;
        // }, Math.floor(Math.random() * 3000));
        
        mediaLoaded = true;
    }
</script>

<div class="{holderClass}" class:hidden={mediaLoaded}></div>

{#if isWebm(src)}
    <video 
        class={mediaClass} 
        class:hidden={!mediaLoaded}
        src={src} 
        autoplay loop muted playsinline 
        on:loadeddata={onMediaLoad}
    ></video>
{:else if src != ""}
    <img 
        class={mediaClass} 
        class:hidden={!mediaLoaded}
        src={src} 
        alt={alt} 
        on:load={onMediaLoad}
    />
{/if}