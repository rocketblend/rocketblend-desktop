<script lang="ts">
    import { AddProjectPackage, RemoveProjectPackage } from '$lib/wailsjs/go/application/Driver'
    import { application } from '$lib/wailsjs/go/models';

    import IconAddFill from '~icons/ri/add-fill';
    import IconSubtractFill from '~icons/ri/subtract-fill';

    export let projectId: string | undefined;
    export let packageRef: string;

    export let assigned: boolean = false;
    export let variantFrom: string = 'secondary';
    export let variantTo: string = 'tertiary';
    export let hovered: boolean = false;
    export let rounded: boolean = true;

    async function togglePackage() {
        if (projectId === undefined) {
            return;
        };

        if (assigned) {
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

    async function handleUserInteraction(event: KeyboardEvent | MouseEvent) {
        if (event.type === 'click' || (event instanceof KeyboardEvent && event.key === 'Enter')) {
            await togglePackage();
        }
    }
</script>

<button
    class="btn h-full bg-gradient-to-br variant-gradient-{variantFrom}-{variantTo} {rounded ? 'rounded' : ''} p-1 text-token"
    on:click|stopPropagation={handleUserInteraction}
    on:keydown|stopPropagation={handleUserInteraction}
>
    {#if hovered}
        {#if assigned}
            <IconSubtractFill class="w-5 h-5" />
        {:else}
            <IconAddFill class="w-5 h-5" />
        {/if}
    {/if}
</button>