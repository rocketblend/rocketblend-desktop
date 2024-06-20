<script lang="ts">
    import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
    import type { ModalSettings, ToastSettings } from "@skeletonlabs/skeleton";

    import { CreateProject } from "$lib/wailsjs/go/application/Driver";
    import { application } from "$lib/wailsjs/go/models"
			
    const modalStore = getModalStore();
    const toastStore = getToastStore();

    export let disabled = false;

    async function handleCreateProject() {
        new Promise<string>((resolve) => {
            const modal: ModalSettings = {
                type: "prompt",
                title: "New Project",
                body: "Pick a name for your new project:",
                value: "New Project",
                buttonTextSubmit: "Create",
                valueAttr: { type: "text", minlength: 3, maxlength: 64, required: true },
                response: (r: string) => {
                    resolve(r);
                }
            };
            modalStore.trigger(modal);
        }).then(async (name: string) => {
            if (!name) {
                return;
            }

            CreateProject(application.CreateProjectOpts.createFrom({
                name: name
            })).then(() => {
                toastStore.trigger({
                    message: "Creating project...",
                });
            });
        });
    }
</script>

<button type="button" class="btn btn-sm variant-filled px-6 font-medium" on:click={handleCreateProject} disabled={disabled}>
    New project
</button>