<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    import { pack } from '$lib/wailsjs/go/models';

    import IconDownload2Fill from '~icons/ri/file-download-fill';
    import IconStopFill from '~icons/ri/stop-mini-fill';
    import IconMoreFill from '~icons/ri/delete-bin-7-fill';
    import IconPlayFill from '~icons/ri/play-fill';
    import IconErrorFill from '~icons/ri/error-warning-fill';

    export let state: pack.PackageState = pack.PackageState.AVAILABLE;
    export let variantFrom: string = 'secondary';
    export let variantTo: string = 'tertiary';
    export let open: boolean = false;
    export let rounded: boolean = true;

    const dispatch = createEventDispatcher();

    function handleAction() {
        switch (state) {
            case pack.PackageState.AVAILABLE:
                dispatch('download');
                break;
            case pack.PackageState.DOWNLOADING:
                dispatch('cancel');
                break;
            case pack.PackageState.INSTALLED:
                dispatch('delete');
                break;
            case pack.PackageState.ERROR:
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
    {#if open}
        {#if state === pack.PackageState.AVAILABLE}
            <IconDownload2Fill />
        {:else if state === pack.PackageState.DOWNLOADING}
            <IconStopFill />
        {:else if state === pack.PackageState.STOPPED}
            <IconPlayFill />
        {:else if state === pack.PackageState.INSTALLED}
            <IconMoreFill />
        {:else if state === pack.PackageState.ERROR}
            <IconErrorFill /> <!-- Example error state handling -->
        {/if}
    {/if}
</div>