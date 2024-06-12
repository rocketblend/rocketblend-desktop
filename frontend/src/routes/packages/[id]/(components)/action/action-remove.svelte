<script lang="ts">
    import { getToastStore, getModalStore, type ToastSettings, type ModalSettings } from '@skeletonlabs/skeleton';

    import { UninstallPackage } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();
    const modalStore = getModalStore();


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

        new Promise<boolean>((resolve) => {
            const modal: ModalSettings = {
                type: "confirm",
                title: "Please Confirm",
                body: "Are you sure you wish to delete the downloaded content for this package?",
                response: (r: boolean) => {
                    resolve(r);
                }
            };
            modalStore.trigger(modal);
        }).then(async (remove: boolean) => {
            if (!remove) {
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

            disabled = false;
        });
    }

    $: displayText = disabled ? disabledText : cancel ? cancelText : deletetext;
</script>

<AlertAction on:click={remove} disabled={disabled}>
    {displayText}
</AlertAction>

