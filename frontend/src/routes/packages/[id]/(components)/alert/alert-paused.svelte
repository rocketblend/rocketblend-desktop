<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

    
    import IconPauseFill from '~icons/ri/pause-fill';

    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
    import { DownloadLabel } from '$lib/components/ui/download';
    import { ActionDownload } from '../action';

    export let packageId: string;
    export let progress: types.Progress | undefined;
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconPauseFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Paused"/>
    </svelte:fragment>
    <AlertDescription message="Package is currently incomplete, resume to finish downloading"/>
    {#if progress}
        <DownloadLabel currentBytes={progress.currentBytes} totalBytes={progress.totalBytes}/>
    {/if}
    <svelte:fragment slot="actions">
        <AlertAction text="Cancel" disabled/>
        <ActionDownload packageId={packageId} text="Resume"/>
    </svelte:fragment>
</Alert>