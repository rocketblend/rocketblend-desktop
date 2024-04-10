<script lang="ts">
    import IconDownloadFill from '~icons/ri/download-2-fill'
    import IconCheckFill from '~icons/ri/check-fill';
    import IconPauseFill from '~icons/ri/pause-fill';
    import IconWarningFill from '~icons/ri/file-damage-fill'

    import { ProgressRadial } from '@skeletonlabs/skeleton';

    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction
    } from '$lib/components/ui/alert';
    import { pack } from '$lib/wailsjs/go/models';

    export let state: pack.PackageState = pack.PackageState.AVAILABLE;

    function cycleState() {
        state = (state + 1) % 5;
    }
</script>

{#if state === pack.PackageState.AVAILABLE}
    <Alert>
        <svelte:fragment slot="icon">
            <IconDownloadFill class="text-2xl"/>
        </svelte:fragment>
        <svelte:fragment slot="title">
            <AlertTitle title="Available"/>
        </svelte:fragment>
        <AlertDescription message="Package is available to be downloaded."/>
        <svelte:fragment slot="actions">
            <AlertAction text="Download" on:click={cycleState}/>
        </svelte:fragment>
    </Alert>
{:else if state === pack.PackageState.DOWNLOADING}
    <Alert>
        <svelte:fragment slot="icon">
            <ProgressRadial width="w-6" stroke={50} strokeLinecap="square"/>
        </svelte:fragment>
        <svelte:fragment slot="title">
            <AlertTitle title="Downloading"/>
        </svelte:fragment>
        <AlertDescription message="Package is currently downloading."/>
        <svelte:fragment slot="actions">
            <AlertAction text="Pause" on:click={cycleState}/>
        </svelte:fragment>
    </Alert>
{:else if state === pack.PackageState.CANCELLED}
    <Alert>
        <svelte:fragment slot="icon">
            <IconPauseFill class="text-2xl"/>
        </svelte:fragment>
        <svelte:fragment slot="title">
            <AlertTitle title="Paused"/>
        </svelte:fragment>
        <AlertDescription message="Download of package has been paused."/>
        <svelte:fragment slot="actions">
            <AlertAction text="Cancel"/>
            <AlertAction text="Resume" on:click={cycleState}/>
        </svelte:fragment>
    </Alert>
{:else if state === pack.PackageState.INSTALLED}
    <Alert>
        <svelte:fragment slot="icon">
            <IconCheckFill class="text-2xl"/>
        </svelte:fragment>
        <svelte:fragment slot="title">
            <AlertTitle title="Ready"/>
        </svelte:fragment>
        <AlertDescription message="Package is ready to be used."/>
        <svelte:fragment slot="actions">
            <AlertAction text="Delete" on:click={cycleState}/>
            <AlertAction text="View Location"/>
        </svelte:fragment>
    </Alert>
{:else}
    <Alert>
        <svelte:fragment slot="icon">
            <IconWarningFill class="text-2xl"/>
        </svelte:fragment>
        <svelte:fragment slot="title">
            <AlertTitle title="Error"/>
        </svelte:fragment>
        <AlertDescription message="An error occurred while downloading the package."/>
        <svelte:fragment slot="actions">
            <AlertAction text="Retry" on:click={cycleState}/>
        </svelte:fragment>
    </Alert>
{/if}