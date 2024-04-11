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
        <AlertTitle title="Not Ready"/>
    </svelte:fragment>
    <AlertDescription message="Package is current enabled on the selected project, but is not downloaded and installed ready for use. See status below."/>
    <svelte:fragment slot="actions">
        <AlertAction text="Disable" variant="glass-surface" on:click/>
    </svelte:fragment>
</Alert>