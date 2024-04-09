<script lang="ts">
    import { onMount } from 'svelte';
    import type { PageData } from './$types';

    import { type ToastSettings, getToastStore  } from '@skeletonlabs/skeleton';

    import type { application } from '$lib/wailsjs/go/models';
	import { UpdateProject } from '$lib/wailsjs/go/application/Driver';

    import { t } from '$lib/translations/translations';
    import { getSelectedProjectStore } from '$lib/stores';
    import { formatDateTime, resourcePath } from '$lib/components/utils';
    import { Media } from '$lib/components/ui/media';
    import { InputInline } from '$lib/components/ui/input';

    import IconEditFill from '~icons/ri/edit-fill';
	import { invalidate } from '$app/navigation';

    const selectedProjectStore = getSelectedProjectStore();
    const toastStore = getToastStore();

    export let data: PageData;

    async function updateProject() {
        const request: application.UpdateProjectOpts = {
            id: data.project.id,
            name: data.project.name || "",
        };

        await UpdateProject(request)
            .then(() => {
                const updateProjectSuccessToast: ToastSettings = {
                    message: t.get('home.toast.saving.save'),
                    timeout: 1000
                };
                toastStore.trigger(updateProjectSuccessToast);
            })
            .catch((error) => {
                const updateProjectErrorToast: ToastSettings = {
                    message: t.get('home.toast.saving.error'),
                    background: "variant-filled-error",
                    timeout: 3000
                };
                toastStore.trigger(updateProjectErrorToast);
                console.error(error);
            });
    }

    function getDependenciesDisplay(): string {
        const dependencies = data.project.addons || [];
        const buildDependencyCount = data.project.build ? 1 : 0;
        return t.get('home.project.tag.dependency', { number: dependencies.length + buildDependencyCount });
    }

    function handleChange(event: CustomEvent) {
        updateProject();
    }

    function setSelectedProject() {
        selectedProjectStore.set([data.project.id.toString()]);
        invalidate("app:layout");
    }

    onMount(() => {
        setSelectedProject();
    });
</script>

<main class="space-y-4"> 
    <div class="flex gap-4 items-end">
        <div>
            <Media src={resourcePath(data.project.thumbnailPath)} alt="" />
        </div>
        <div class="space-y-2">
            <InputInline bind:value={data.project.name} labelClasses="h2 font-bold items-baseline" inputClasses="input" on:change={handleChange}>
                <IconEditFill class="text-sm text-surface-600-300-token"/>
            </InputInline>
            <div class="flex flex-wrap text-sm text-surface-800-100-token gap-1">
                <div class="badge variant-ghost rounded">{data.project.path}</div>
                <div class="badge variant-ghost rounded">{data.project.id}</div>
                <div class="badge variant-ghost rounded">{data.project.fileName}</div>
                <div class="badge variant-ghost rounded">{data.project.build}</div>
                {#each data.project.tags || [] as tag}
                    <div class="badge variant-ghost-primary rounded">{tag}</div>
                {/each}
                <div class="badge variant-ghost-secondary rounded">{getDependenciesDisplay()}</div>
                <div class="badge variant-ghost rounded">{formatDateTime(data.project.updatedAt)}</div>
            </div>
        </div>
    </div>
    <!-- <hr>
    <InlineInput type="textarea" placeholder="Add description..."/> -->
    <hr>
    <div class="grid grid-cols-4 gap-4">
        <Media height="80" width="full" src={resourcePath(data.project.splashPath)} />
      </div>
</main>