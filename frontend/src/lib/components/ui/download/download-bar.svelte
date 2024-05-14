<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

	import ProgressBar from '$lib/components/ui/progress/progress-bar.svelte';
	import FormatBytes from '$lib/components/ui/format/format-bytes.svelte';

    import DownloadLabel from './download-label.svelte';

    export let progress: types.Progress;
    export let size: "sm" | "md" | "lg" = "sm";
</script>

{#if progress.currentBytes != progress.totalBytes}
    <p class="text-{size}">Retrieving package data...</p>
    <ProgressBar value={progress.currentBytes} max={progress.totalBytes} >
        <DownloadLabel currentBytes={progress.currentBytes} totalBytes={progress.totalBytes}/>
    </ProgressBar>
    <FormatBytes value={progress.bytesPerSecond} suffix="/s"/>
{:else}
    <p class="text-{size}">Wrapping things up...</p>
{/if}