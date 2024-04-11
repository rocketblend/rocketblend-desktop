<script lang="ts">
    import IconCheckFill from '~icons/ri/check-fill';
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconCheckFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Ready"/>
    </svelte:fragment>
    <AlertDescription message="Package is ready to be used."/>
    <svelte:fragment slot="actions">
        <AlertAction text="Delete" on:click/>
        <AlertAction text="View Location"/>
    </svelte:fragment>
</Alert>