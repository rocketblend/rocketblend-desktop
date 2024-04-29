<script lang="ts">
    import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';

    import IconDownloadFill from '~icons/ri/download-2-fill';

    import { InstallPackage } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();

    export let packageId: string

    let disabled = false;

    function downloadPackage() {
        if (disabled) {
            return;
        }

        disabled = true;
        const opts = application.InstallPackageOpts.createFrom({ id: packageId });

        InstallPackage(opts).then((result) => {
            // TODO: Handle these types of notifiations globally via events from go.
            const downloadPackageToast: ToastSettings = {
                message: `Downloading Package: ${packageId}`,
            };

            toastStore.trigger(downloadPackageToast);
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error starting package download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });

        disabled = false;
    }
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
        <AlertAction text="Download" disabled={disabled} on:click={downloadPackage}/>
    </svelte:fragment>
</Alert>