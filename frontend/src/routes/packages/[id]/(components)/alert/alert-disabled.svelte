<script lang="ts">
    import IconCloseFill from '~icons/ri/close-fill';
    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction
    } from '$lib/components/ui/alert';
</script>
  
<Alert>
    <svelte:fragment slot="icon">
        <IconCloseFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Disabled
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        Package is currently disabled on the selected project.
    </AlertDescription>
    <svelte:fragment slot="actions">
        <AlertAction on:click>
            Enabled
        </AlertAction>
    </svelte:fragment>
</Alert>