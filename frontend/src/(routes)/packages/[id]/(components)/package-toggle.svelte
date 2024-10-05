<script lang="ts">
    import { AddProjectPackage, RemoveProjectPackage } from '$lib/wailsjs/go/application/Driver'
    import { enums, application, } from '$lib/wailsjs/go/models';

    import { AlertEnabled, AlertDisabled, AlertNotReady } from './alert';

    export let projectId: string | undefined;
    export let projectBuildRef: string | undefined;
    export let projectAddonRefs: string[] | undefined;
    export let packageRef: string;
    export let packageState: enums.PackageState;

    async function togglePackage() {
        if (projectId === undefined) {
            return;
        };

        if (isActive) {
            await RemoveProjectPackage(application.RemoveProjectPackageOpts.createFrom({
                id: projectId,
                reference: packageRef,
            }));
            
            return;
        }

        await AddProjectPackage(application.AddProjectPackageOpts.createFrom({
            id: projectId,
            reference: packageRef,
        }));
    }

    $: isActiveBuild = projectBuildRef=== packageRef;
    $: isActiveAddon = !!projectAddonRefs?.some(ref => ref === packageRef);
    $: isActive = isActiveBuild || isActiveAddon;
</script>

{#if isActive}
    {#if packageState === enums.PackageState.INSTALLED}
        <AlertEnabled on:click={togglePackage}/>
    {:else}
        <AlertNotReady on:click={togglePackage}/>
    {/if}
{:else}
    <AlertDisabled on:click={togglePackage}/>
{/if}