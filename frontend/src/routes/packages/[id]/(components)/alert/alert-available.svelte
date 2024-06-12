<script lang="ts">
    import IconDownloadFill from '~icons/ri/download-2-fill';

    import { Alert, AlertTitle, AlertDescription } from '$lib/components/ui/alert';
    import { ActionDownload } from '../action';

    export let packageId: string
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconDownloadFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle>
            Available
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription message="Package is available to be downloaded."/>
    <svelte:fragment slot="actions">
        <ActionDownload packageId={packageId}/>
    </svelte:fragment>
</Alert>