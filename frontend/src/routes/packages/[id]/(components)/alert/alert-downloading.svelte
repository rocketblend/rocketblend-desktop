<script lang="ts">
    import { ProgressRadial } from '@skeletonlabs/skeleton';
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
</script>

<Alert>
    <svelte:fragment slot="icon">
        <ProgressRadial width="w-6" stroke={50} strokeLinecap="square"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Downloading"/>
    </svelte:fragment>
    <AlertDescription message="Package is currently downloading."/>
    <svelte:fragment slot="actions">
        <AlertAction text="Pause" on:click/>
    </svelte:fragment>
</Alert>