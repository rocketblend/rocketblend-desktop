<script lang="ts">
    import IconInfoFill from '~icons/ri/check-double-fill';
    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction
    } from '$lib/components/ui/alert';
</script>
  
<Alert variant="ghost-primary">
    <svelte:fragment slot="icon">
        <IconInfoFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Enabled
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        Package is currently enabled on the selected project.
    </AlertDescription>
    <svelte:fragment slot="actions">
        <AlertAction variant="glass-surface" on:click>
            Disabled
        </AlertAction>
    </svelte:fragment>
</Alert>