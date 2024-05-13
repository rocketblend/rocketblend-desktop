<script lang="ts">
    import { ProgressRadial } from '@skeletonlabs/skeleton';
    import type { types } from '$lib/wailsjs/go/models';

    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
	import ProgressBar from '$lib/components/ui/progress/progress-bar.svelte';
	import FormatBytes from '$lib/components/ui/format/format-bytes.svelte';

    export let progress: types.Progress | undefined;
</script>

<Alert>
    <svelte:fragment slot="icon">
        <ProgressRadial width="w-6" stroke={50} strokeLinecap="square"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Downloading"/>
    </svelte:fragment>
    <AlertDescription message="Package is currently downloading."/>
    {#if progress}
        <ProgressBar value={progress.currentBytes} max={progress.totalBytes} >
            <div class="flex items-center gap-2 whitespace-nowrap">
                <FormatBytes value={progress.currentBytes}/>
                <div> / </div>
                <FormatBytes value={progress.totalBytes}/>
            </div>
        </ProgressBar>
        <FormatBytes value={progress.bytesPerSecond} suffix="/s"/>
    {/if}
    <svelte:fragment slot="actions">
        <AlertAction text="Pause" disabled/>
    </svelte:fragment>
</Alert>