<script lang="ts">
    import { ProgressRadial } from '@skeletonlabs/skeleton';
    import type { types } from '$lib/wailsjs/go/models';

    import { Alert, AlertTitle, AlertAction } from '$lib/components/ui/alert';
	import { DownloadBar } from '$lib/components/ui/download';

    export let progress: types.Progress | undefined;
</script>

<Alert>
    <svelte:fragment slot="icon">
        <ProgressRadial width="w-6" stroke={50} strokeLinecap="square"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Downloading"/>
    </svelte:fragment>
    {#if progress}
        <DownloadBar
            currentBytes={progress.currentBytes}
            totalBytes={progress.totalBytes}
            bytesPerSecond={progress.bytesPerSecond}
        />
    {/if}
    <svelte:fragment slot="actions">
        {#if progress && progress.currentBytes != progress.totalBytes }
            <AlertAction text="Pause" disabled/>
        {/if}
    </svelte:fragment>
</Alert>