<script lang="ts">
    import IconCloseFill from '~icons/ri/close-fill';
    import IconInfoFill from '~icons/ri/check-double-fill';
    import IconWarningFill from '~icons/ri/error-warning-fill'

    import { pack } from '$lib/wailsjs/go/models';

    export let active = false;
    export let state: pack.PackageState = pack.PackageState.AVAILABLE;

    function toggle() {
        active = !active;
    }
</script>

{#if active}
    {#if state === pack.PackageState.INSTALLED}
        <aside class="alert variant-ghost-primary">
            <IconInfoFill class="text-2xl"/>
            <div class="alert-message ">
                <h2 class="font-bold h6">Enabled</h2>
                <p class="text-sm">Package is current enabled on the selected project</p>
            </div>
            <div class="alert-actions">
                <button class="btn btn-sm variant-glass-surface font-medium" on:click={toggle}>Disable</button>
            </div>
        </aside>
    {:else}
        <aside class="alert variant-ghost-warning">
            <IconWarningFill class="text-2xl"/>
            <div class="alert-message ">
                <h2 class="font-bold h6">Not Ready</h2>
                <p class="text-sm">Package is current enabled on the selected project, but is not downloaded and installed ready for use. See status below.</p>
            </div>
            <div class="alert-actions">
                <button class="btn btn-sm variant-glass-surface font-medium" on:click={toggle}>Disable</button>
            </div>
        </aside>
    {/if}
{:else}
    <aside class="alert variant-ghost-surface">
        <IconCloseFill class="text-2xl"/>
        <div class="alert-message ">
            <h2 class="font-bold h6">Disabled</h2>
            <p class="text-sm">Package is currently disabled on the selected project</p>
        </div>
        <div class="alert-actions">
            <button class="btn btn-sm variant-filled-surface font-medium" on:click={toggle}>Enable</button>
        </div>
    </aside>
{/if}