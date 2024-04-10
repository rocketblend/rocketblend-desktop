<script lang="ts">
    import IconCloseFill from '~icons/ri/close-fill';
    import IconInfoFill from '~icons/ri/check-double-fill';
    import IconWarningFill from '~icons/ri/error-warning-fill'

    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction
    } from '$lib/components/ui/alert';
    import { pack } from '$lib/wailsjs/go/models';

    export let active = false;
    export let state: pack.PackageState = pack.PackageState.AVAILABLE;
</script>

{#if active}
    {#if state === pack.PackageState.INSTALLED}
        <Alert variant="ghost-primary">
            <svelte:fragment slot="icon">
                <IconInfoFill class="text-2xl"/>
            </svelte:fragment>
            <svelte:fragment slot="title">
                <AlertTitle title="Enabled"/>
            </svelte:fragment>
            <AlertDescription message="Package is current enabled on the selected project."/>
            <svelte:fragment slot="actions">
                <AlertAction text="Disable" variant="glass-surface" on:click/>
            </svelte:fragment>
        </Alert>
    {:else}
        <Alert variant="ghost-warning">
            <svelte:fragment slot="icon">
                <IconWarningFill class="text-2xl"/>
            </svelte:fragment>
            <svelte:fragment slot="title">
                <AlertTitle title="Not Ready"/>
            </svelte:fragment>
            <AlertDescription message="Package is current enabled on the selected project, but is not downloaded and installed ready for use. See status below."/>
            <svelte:fragment slot="actions">
                <AlertAction text="Disable" variant="glass-surface" on:click/>
            </svelte:fragment>
        </Alert>
    {/if}
{:else}
    <Alert>
        <svelte:fragment slot="icon">
            <IconCloseFill class="text-2xl"/>
        </svelte:fragment>
        <svelte:fragment slot="title">
            <AlertTitle title="Disabled"/>
        </svelte:fragment>
        <AlertDescription message="Package is currently disabled on the selected project"/>
        <svelte:fragment slot="actions">
            <AlertAction text="Enable" on:click/>
        </svelte:fragment>
    </Alert>
{/if}