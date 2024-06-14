<script lang="ts">
    import { createEventDispatcher } from 'svelte';
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
</script>

<button
    type="button"
    on:click={handleClick}
    on:keydown={(e) => e.key === 'Enter' && handleClick()}
    class={`media-button ${className} ${hover ? 'img-hover' : ''}`}
>
    <img src={src} alt={alt} loading={loading} class="media-image" />
</button>

<style>
    .media-button {
        border: none;
        background: none;
        padding: 0;
        cursor: pointer;
        width: 100%;
        display: block;
    }
    
    .media-button .media-image {
        width: 100%;
    }

    .img-hover {
        opacity: 0.9;
        transition: all 0.2s;
    }

    .img-hover:hover,
    .media-button:focus {
        opacity: 1;
        transform: scale(1.05);
    }
</style>
