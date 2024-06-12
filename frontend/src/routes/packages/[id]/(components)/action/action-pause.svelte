<script lang="ts">
    import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';

    import { CancelOperation } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();

    export let downloadId: string | undefined;

    let text: string = "Pause";
    let disabledText: string = "Pausing...";

    let disabled = false;

    function pause() {
        if (disabled || !downloadId) {
            return;
        }

        disabled = true;
        const opts = application.CancelOperationOpts.createFrom({ id: downloadId });

        CancelOperation(opts).then(() => {
            return;
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error pausing download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });
    }

    $: displayText = disabled ? disabledText : text;
</script>

<AlertAction on:click={pause} disabled={disabled || !downloadId}>
    {displayText}
</AlertAction>

