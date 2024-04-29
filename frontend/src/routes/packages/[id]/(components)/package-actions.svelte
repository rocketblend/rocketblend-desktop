<script lang="ts">
    import { enums } from '$lib/wailsjs/go/models';

    import {
        AlertAvailable,
        AlertDownloading,
        AlertPaused,
        AlertInstalled,
        AlertError
    } from './alert';

    export let packageId: string;
    export let state: enums.PackageState = enums.PackageState.AVAILABLE;
</script>

{#if state === enums.PackageState.AVAILABLE}
    <AlertAvailable packageId={packageId}/>
{:else if state === enums.PackageState.DOWNLOADING}
    <AlertDownloading />
{:else if state === enums.PackageState.CANCELLED}
    <AlertPaused />
{:else if state === enums.PackageState.INSTALLED}
    <AlertInstalled />
{:else}
    <AlertError />
{/if}