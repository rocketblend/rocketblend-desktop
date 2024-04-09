<script lang="ts">
    import IconDownloadFill from '~icons/ri/download-2-fill'
    import IconCheckFill from '~icons/ri/check-fill';
    import IconPauseFill from '~icons/ri/pause-fill';
    import IconWarningFill from '~icons/ri/file-damage-fill'

    import { ProgressRadial } from '@skeletonlabs/skeleton';

    import { pack } from '$lib/wailsjs/go/models';

    export let state: pack.PackageState = pack.PackageState.AVAILABLE;

    function cycleState() {
        state = (state + 1) % 5;
    }

</script>

{#if state === pack.PackageState.AVAILABLE}
    <aside class="alert variant-ghost-surface">
        <IconDownloadFill class="text-2xl"/>
        <div class="alert-message ">
            <h2 class="font-bold h6">Available</h2>
            <p class="text-sm">Package is available to be downloaded.</p>
        </div>
        <div class="alert-actions">
            <button class="btn btn-sm variant-filled-surface font-medium" on:click={cycleState}>Download</button>
        </div>
    </aside>
{:else if state === pack.PackageState.DOWNLOADING}
    <aside class="alert variant-ghost-surface">
        <ProgressRadial width="w-6" stroke={50} strokeLinecap="square"/>
        <div class="alert-message ">
            <h2 class="font-bold h6">Downloading</h2>
            <p class="text-sm">Package is currently downloading.</p>
        </div>
        <div class="alert-actions">
            <button class="btn btn-sm variant-filled-surface font-medium" on:click={cycleState}>Pause</button>
        </div>
    </aside>
{:else if state === pack.PackageState.CANCELLED}
    <aside class="alert variant-ghost-surface">
        <IconPauseFill class="text-2xl"/>
        <div class="alert-message ">
            <h2 class="font-bold h6">Paused</h2>
            <p class="text-sm">Download of package has been paused.</p>
        </div>
        <div class="alert-actions">
            <button class="btn btn-sm variant-filled-surface font-medium">Cancel</button>
            <button class="btn btn-sm variant-filled-surface font-medium" on:click={cycleState}>Resume</button>
        </div>
    </aside>
{:else if state === pack.PackageState.INSTALLED}
    <aside class="alert variant-ghost-surface">
        <IconCheckFill class="text-2xl"/>
        <div class="alert-message ">
            <h2 class="font-bold h6">Ready</h2>
            <p class="text-sm">Package is ready to be used.</p>
        </div>
        <div class="alert-actions">
            <button class="btn btn-sm variant-filled-surface font-medium"  on:click={cycleState}>Delete</button>
            <button class="btn btn-sm variant-filled-surface font-medium">View Location</button>
        </div>
    </aside>
{:else}
    <aside class="alert variant-ghost-surface">
        <IconWarningFill class="text-2xl"/>
        <div class="alert-message ">
            <h2 class="font-bold h6">Error</h2>
            <p class="text-sm">An error occurred while downloading the package.</p>
        </div>
        <div class="alert-actions">
            <button class="btn btn-sm variant-filled-surface font-medium" on:click={cycleState}>Retry</button>
        </div>
    </aside>
{/if}