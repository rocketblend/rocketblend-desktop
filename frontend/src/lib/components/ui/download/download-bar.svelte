<script lang="ts">
    import ProgressBar from '$lib/components/ui/progress/progress-bar.svelte';
    import FormatBytes from '$lib/components/ui/format/format-bytes.svelte';

    import DownloadLabel from './download-label.svelte';

    export let currentBytes: number;
    export let totalBytes: number;
    export let bytesPerSecond: number;
    export let size: "sm" | "md" | "lg" = "sm";
</script>

{#if currentBytes != totalBytes}
    <p class="text-{size}">Retrieving package data...</p>
    <ProgressBar value={currentBytes} max={totalBytes} >
        <DownloadLabel currentBytes={currentBytes} totalBytes={totalBytes}/>
    </ProgressBar>
    <FormatBytes value={bytesPerSecond} suffix="/s"/>
{:else}
    <p class="text-{size}">Wrapping things up...</p>
{/if}