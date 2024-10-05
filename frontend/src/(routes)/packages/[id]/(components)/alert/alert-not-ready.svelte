<script lang="ts">
    import IconWarningFill from '~icons/ri/error-warning-fill'
    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction
    } from '$lib/components/ui/alert';
</script>

<Alert variant="ghost-warning">
    <svelte:fragment slot="icon">
        <IconWarningFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle>
            Not Ready
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        Package is current enabled on the selected project, but is not downloaded and installed ready for use. See status below.
    </AlertDescription>
    <svelte:fragment slot="actions">
        <AlertAction variant="glass-surface" on:click>
            Disabled 
        </AlertAction>
    </svelte:fragment>
</Alert>