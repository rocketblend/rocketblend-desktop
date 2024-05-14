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

    function download() {
        if (disabled) {
            return;
        }

        disabled = true;
        const opts = application.InstallPackageOpts.createFrom({ id: packageId });

        InstallPackage(opts).then(() => {
            const downloadPackageToast: ToastSettings = {
                message: `Download starting...`,
                timeout: 3000,
            };

            toastStore.trigger(downloadPackageToast);
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error starting download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });
    }

    $: displayText = disabled ? disabledText : text;
</script>

<AlertAction bind:text={displayText} on:click={download} disabled={disabled}/>
