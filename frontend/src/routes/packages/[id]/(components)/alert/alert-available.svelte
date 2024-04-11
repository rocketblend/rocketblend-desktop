<script lang="ts">
    import IconDownloadFill from '~icons/ri/download-2-fill';
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconDownloadFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Available"/>
    </svelte:fragment>
    <AlertDescription message="Package is available to be downloaded."/>
    <svelte:fragment slot="actions">
        <AlertAction text="Download" on:click/>
    </svelte:fragment>
</Alert>