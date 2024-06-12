<script lang="ts">
    import { ProgressRadial } from '@skeletonlabs/skeleton';
    import type { types } from '$lib/wailsjs/go/models';

    import { Alert, AlertTitle, AlertDescription } from '$lib/components/ui/alert';
	import { DownloadBar } from '$lib/components/ui/download';
    import { ActionPause } from '../action';

    export let progress: types.Progress | undefined;
    export let downloadId: string | undefined;

    let title = "Downloading"

    $: displayTitle = downloadId ? title : `${title} (External)`;
</script>

<Alert>
    <svelte:fragment slot="icon">
        <ProgressRadial width="w-6" stroke={50} strokeLinecap="square"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            {displayTitle}
        </AlertTitle>
    </svelte:fragment>
    {#if !downloadId}
        <AlertDescription message="Package is being downloaded externally, so cannot be paused here."/>
    {/if}
    {#if progress}
        <DownloadBar
            currentBytes={progress.currentBytes}
            totalBytes={progress.totalBytes}
            bytesPerSecond={progress.bytesPerSecond}
        />
    {/if}
    <svelte:fragment slot="actions">
        {#if progress && progress.currentBytes != progress.totalBytes }
            <ActionPause downloadId={downloadId}/>
        {/if}
    </svelte:fragment>
</Alert>