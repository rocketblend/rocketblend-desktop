<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { PackageState } from './types';

    import IconDownload2Fill from '~icons/ri/file-download-fill';
    import IconStopFill from '~icons/ri/stop-mini-fill';
    import IconMoreFill from '~icons/ri/delete-bin-7-fill';
    import IconError from '~icons/ri/error-warning-fill';

    export let state: PackageState = PackageState.Available;
    export let variantFrom: string = 'secondary';
    export let variantTo: string = 'tertiary';
    export let isOpen: boolean = false;
    export let rounded: boolean = true;

    const dispatch = createEventDispatcher();

    function handleAction() {
        switch (state) {
            case PackageState.Available:
                dispatch('download');
                break;
            case PackageState.Downloading:
                dispatch('cancel');
                break;
            case PackageState.Installed:
                dispatch('delete');
                break;
            case PackageState.Error:
                // Handle error state if needed
                break;
            // Add additional cases as necessary
        }
    }

    function handleUserInteraction(event: KeyboardEvent | MouseEvent) {
        if (event.type === 'click' || (event instanceof KeyboardEvent && event.key === 'Enter')) {
            handleAction();
        }
    }
</script>

<div
    class="flex items-center h-full bg-gradient-to-br variant-gradient-{variantFrom}-{variantTo} {rounded ? 'rounded' : ''} p-1 text-token"
    on:click|stopPropagation={handleUserInteraction}
    on:keydown|stopPropagation={handleUserInteraction}
    role="button"
    tabindex="0"
>
    {#if isOpen}
        {#if state === PackageState.Available}
            <IconDownload2Fill />
        {:else if state === PackageState.Downloading}
            <IconStopFill />
        {:else if state === PackageState.Installed}
            <IconMoreFill />
        {:else if state === PackageState.Error}
            <IconError /> <!-- Example error state handling -->
        {/if}
    {/if}
</div>