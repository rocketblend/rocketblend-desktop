<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction,
    } from '$lib/components/ui/alert';

    export let id: string = "";
    export let reference: string;
</script>
  
<Alert variant="ghost-warning">
    <svelte:fragment slot="title">
        <AlertTitle >
            Missing dependency
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        The dependency <span class="font-semibold">"{reference}"</span> is missing on your machine. Please install it before continuing.
    </AlertDescription>
    <svelte:fragment slot="actions">
        {#if id}
            <AlertAction>
                View
            </AlertAction>
        {/if}
    </svelte:fragment>
</Alert>