<script lang="ts">
    import IconErrorWarningFill from '~icons/ri/error-warning-fill'

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
    <svelte:fragment slot="icon">
        <IconErrorWarningFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Missing dependency
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        The dependency <span class="font-semibold">"{reference}"</span> is missing on your computer. Please install it before continuing.
    </AlertDescription>
    <svelte:fragment slot="actions">
        {#if id}
            <AlertAction href="/packages/{id}" variant="glass-surface">
                View
            </AlertAction>
        {/if}
    </svelte:fragment>
</Alert>