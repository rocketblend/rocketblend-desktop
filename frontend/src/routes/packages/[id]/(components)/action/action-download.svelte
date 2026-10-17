<script lang="ts">
    import { getToastStore, getModalStore, type ToastSettings, type ModalSettings } from '@skeletonlabs/skeleton';

    import { InstallPackage } from '$lib/wailsjs/go/application/Package';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    const toastStore = getToastStore();
    const modalStore = getModalStore();

    export let packageId: string;
    export let resume: boolean = false;
    
    let startText: string = "Download";
    let resumeText: string = "Resume";
    let disabledText: string = "Starting...";

    let disabled = false;

    function confirmDomain(domain: string): Promise<boolean> {
        return new Promise<boolean>((resolve) => {
            const modal: ModalSettings = {
                type: "confirm",
                title: "Untrusted Download",
                body: `This package downloads from <strong>${domain}</strong>, which isn't a trusted domain. Only continue if you trust this source.`,
                buttonTextConfirm: "Download",
                response: (r: boolean) => {
                    resolve(r);
                }
            };
            modalStore.trigger(modal);
        });
    }

    async function install(allowUntrustedDomain: boolean) {
        const opts = application.InstallPackageOpts.createFrom({ id: packageId, allowUntrustedDomain });
        const result = await InstallPackage(opts);

        if (result.requiresConfirmation && !allowUntrustedDomain) {
            if (!(await confirmDomain(result.domain ?? "an unknown domain"))) {
                disabled = false;
                return;
            }

            await install(true);
        }
    }

    function download() {
        if (disabled) {
            return;
        }

        disabled = true;

        install(false).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error starting download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
            disabled = false;
        });
    }

    $: displayText = disabled ? disabledText : resume ? resumeText : startText;
</script>

<AlertAction on:click={download} disabled={disabled}>
    {displayText}
</AlertAction>
