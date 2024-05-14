<script lang="ts">
    import { enums, type types } from '$lib/wailsjs/go/models';

    import {
        AlertAvailable,
        AlertDownloading,
        AlertPaused,
        AlertInstalled,
        AlertError
    } from './alert';

    export let packageId: string;
    export let installationPath: string;
    export let state: enums.PackageState = enums.PackageState.AVAILABLE;
    export let progress: types.Progress | undefined;
</script>

{#if state === enums.PackageState.AVAILABLE}
    <AlertAvailable packageId={packageId}/>
{:else if state === enums.PackageState.DOWNLOADING}
    <AlertDownloading progress={progress} />
{:else if state === enums.PackageState.INCOMPLETE}
    <AlertPaused packageId={packageId} progress={progress} />
{:else if state === enums.PackageState.INSTALLED}
    <AlertInstalled path={installationPath}/>
{:else}
    <AlertError />
{/if}