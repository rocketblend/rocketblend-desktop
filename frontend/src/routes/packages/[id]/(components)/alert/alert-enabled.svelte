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
        <AlertTitle title="Enabled"/>
    </svelte:fragment>
    <AlertDescription message="Package is currently enabled on the selected project."/>
    <svelte:fragment slot="actions">
        <AlertAction text="Disable" variant="glass-surface" on:click/>
    </svelte:fragment>
</Alert>