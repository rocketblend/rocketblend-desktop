<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';

    import { InstallPackage } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();

    export let packageId: string;
    export let text: string = "Download";
    export let disabledText: string = "Starting...";

    let disabled = false;

    const dispatch = createEventDispatcher();

    function download() {
        if (disabled) {
            return;
        }

        disabled = true;
        const opts = application.InstallPackageOpts.createFrom({ id: packageId });

        InstallPackage(opts).then((result) => {
            const downloadPackageToast: ToastSettings = {
                message: `Downloading Package: ${packageId}`,
            };

            toastStore.trigger(downloadPackageToast);
            dispatch('download-started', { packageId });
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error starting package download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });
    }

    $: displayText = disabled ? disabledText : text;
</script>

<AlertAction bind:text={displayText} on:click={download} disabled={disabled}/>
