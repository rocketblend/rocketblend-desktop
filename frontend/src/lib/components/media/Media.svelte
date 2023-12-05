<script lang="ts">
    import { twMerge } from 'tailwind-merge';

    export let src: string;
    export let alt: string = '';
    export let mediaClass: string = 'rounded-container-token';

    $: tagClass = twMerge(mediaClass, $$props.class);

    const isWebm = (path: string) => path.endsWith('.webm');
</script>

{#if isWebm(src)}
    <video {...$$restProps} class={tagClass} src={src} autoplay loop muted playsinline></video>
{:else}
    <img {...$$restProps} class={tagClass} src={src} alt={alt} />
{/if}