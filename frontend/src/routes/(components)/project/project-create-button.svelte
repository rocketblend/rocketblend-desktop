<script lang="ts">
    import { getModalStore } from "@skeletonlabs/skeleton";
    import type { ModalSettings } from "@skeletonlabs/skeleton";

    import { CreateProjectOperation } from "$lib/wailsjs/go/application/Driver";
    import { application } from "$lib/wailsjs/go/models"
			
    const modalStore = getModalStore();

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

            console.log("resolved response:", name);

            const opts = application.CreateProjectOperationOpts.createFrom({ name: name });
            await CreateProjectOperation(opts);

            console.log("Project creation started!")
        });
    }
</script>

<button type="button" class="btn btn-sm variant-filled px-6 font-medium" on:click={handleCreateProject}>
    New project
</button>