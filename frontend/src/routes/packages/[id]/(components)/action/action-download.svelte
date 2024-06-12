<script lang="ts">
    import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';

    import { InstallPackage } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();

    export let packageId: string;
    export let resume: boolean = false;
    
    let startText: string = "Download";
    let resumeText: string = "Resume";
    let disabledText: string = "Starting...";

    let disabled = false;

    function download() {
        if (disabled) {
            return;
        }

        disabled = true;
        const opts = application.InstallPackageOpts.createFrom({ id: packageId });

        InstallPackage(opts).then(() => {
            return;
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error starting download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });
    }

    $: displayText = disabled ? disabledText : resume ? resumeText : startText;
</script>

<AlertAction on:click={download} disabled={disabled}>
    {displayText}
</AlertAction>
