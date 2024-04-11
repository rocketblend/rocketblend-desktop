<script lang="ts">
    import { pack } from '$lib/wailsjs/go/models';

    import {
        AlertAvailable,
        AlertDownloading,
        AlertPaused,
        AlertInstalled,
        AlertError
    } from './alert';

    export let state: pack.PackageState = pack.PackageState.AVAILABLE;
</script>

{#if state === pack.PackageState.AVAILABLE}
    <AlertAvailable />
{:else if state === pack.PackageState.DOWNLOADING}
    <AlertDownloading />
{:else if state === pack.PackageState.CANCELLED}
    <AlertPaused />
{:else if state === pack.PackageState.INSTALLED}
    <AlertInstalled />
{:else}
    <AlertError />
{/if}