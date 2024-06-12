<script lang="ts">
    import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';

    import { UninstallPackage } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();

    export let packageId: string;
    export let cancel: boolean = false;

    let deletetext: string = "Delete";
    let cancelText: string = "Cancel";
    let disabledText: string = "Removing...";

    let disabled = false;

    function remove() {
        if (disabled) {
            return;
        }

        disabled = true;
        const opts = application.UninstallPackageOpts.createFrom({ id: packageId });

        UninstallPackage(opts).then(() => {
            return;
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error removing installation: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });
    }

    $: displayText = disabled ? disabledText : cancel ? cancelText : deletetext;
</script>

<AlertAction on:click={remove} disabled={disabled}>
    {displayText}
</AlertAction>

