<script lang="ts">
    import { getModalStore } from '@skeletonlabs/skeleton';
    import type { ModalSettings } from '@skeletonlabs/skeleton';
			
    const modalStore = getModalStore();

    async function handleCreateProject() {
        new Promise<string>((resolve) => {
            const modal: ModalSettings = {
                type: 'prompt',
                title: 'New Project',
                body: 'Provide the name you want to give to your new project.',
                value: 'New Project',
                valueAttr: { type: 'text', minlength: 3, maxlength: 10, required: true },
                response: (r: string) => {
                    resolve(r);
                }
            };
            modalStore.trigger(modal);
        }).then((r: any) => {
            // Call Create project endpoint
            console.log('resolved response:', r);
        });
    }
</script>

<button type="button" class="btn btn-sm variant-filled px-6 font-medium" on:click={handleCreateProject}>
    New project
</button>