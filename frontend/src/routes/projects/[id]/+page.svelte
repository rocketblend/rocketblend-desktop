<script lang="ts">
    import { onMount } from 'svelte';
    import type { PageData } from './$types';
    import { invalidate } from '$app/navigation';

    import { type ToastSettings, type ModalSettings, getToastStore, getModalStore  } from '@skeletonlabs/skeleton';

    import type { application, types } from '$lib/wailsjs/go/models';
	import { UpdateProject } from '$lib/wailsjs/go/application/Driver';

    import { t } from '$lib/translations/translations';
    import { getSelectedProjectStore } from '$lib/stores';
    import { formatDateTime } from '$lib/components/utils';
    import { InputInline } from '$lib/components/ui/input';
    import { Gallery, Media, type GalleryItem } from '$lib/components/ui/gallery';

    import { AlertEmptyMedia } from './(components)/alert';

    import IconEditFill from '~icons/ri/edit-fill';

    const selectedProjectStore = getSelectedProjectStore();
    const toastStore = getToastStore();
    const modalStore = getModalStore();

    export let data: PageData;

    let dependenciesLabel: string;

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

    function handleChange(event: CustomEvent) {
        updateProject();
    }

    function setSelectedProject() {
        selectedProjectStore.set([data.project.id.toString()]);
        invalidate("app:layout");
    }

    function convertToGalleryItems(projects: types.Media[] = []): GalleryItem[] {
        return projects.map((media) => ({
            value: media.filePath || "",
            src: media.url || "",
            alt: `${media.filePath || ""}`,
            class: "",
        }));
    }

    function handleGalleryClick(event: CustomEvent<{ value: string }>) {
        if (!data.project.media) {
        return;
    }

        const filepath = event.detail.value;
        const index = data.project.media.findIndex((m) => m.filePath === filepath);

        if (index === -1) {
            //console.error('Media item not found');
            return;
        }

        const modal: ModalSettings = {
            type: 'component',
            component: 'modalMediaViewer',
            modalClasses: "h-full",
            meta: {
                media: data.project.media,
                goto: index,
            },
        };

        modalStore.trigger(modal);
    }

    onMount(() => {
        setSelectedProject();
    });

    $: {
        const dependencies = data.project.addons || [];
        const buildDependencyCount = data.project.build ? 1 : 0;
        dependenciesLabel = t.get('home.project.tag.dependency', { number: dependencies.length + buildDependencyCount })
    }

    $: updatedAt = formatDateTime(data.project.updatedAt);
    $: galleryItems = convertToGalleryItems(data.project.media || []);
</script>

<main class="flex flex-col h-full space-y-4"> 
    <div class="flex gap-4 items-end">
        <div>
            <Media src={data.project.thumbnail?.url} height={32} width={32} class="cursor-default" rounded/>
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
                <div class="badge variant-ghost-secondary rounded">{dependenciesLabel}</div>
                <div class="badge variant-ghost rounded">{updatedAt}</div>
            </div>
        </div>
    </div>
    <hr>
    <div class="h-full overflow-auto space-y-4">
        {#if galleryItems.length > 0}
            <Gallery
                gap={15}
                maxColumnWidth={250}
                bind:items={galleryItems}
                on:click={handleGalleryClick}
                loading="eager"
                rounded
            />
        {:else}
            <AlertEmptyMedia folder={data.project.mediaPath}/>
        {/if}
        <p class="text-sm text-surface-600-300-token">Want to set a specific file as either the splash or the thumbnail? Just add <code class="code">splash</code> or <code class="code">thumbnail</code> respectively to the filename.</p>
    </div>
</main>