<script lang="ts">
    import IconPauseFill from '~icons/ri/pause-fill';
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconPauseFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Paused"/>
    </svelte:fragment>
    <AlertDescription message="Download of package has been paused."/>
    <svelte:fragment slot="actions">
        <AlertAction text="Cancel"/>
        <AlertAction text="Resume" disabled/>
    </svelte:fragment>
</Alert>