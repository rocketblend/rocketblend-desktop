<script lang="ts">
    import { onMount } from 'svelte';
    import type { PageData } from './$types';
    import { invalidate } from '$app/navigation';

    import { type ToastSettings, type ModalSettings, getToastStore, getModalStore  } from '@skeletonlabs/skeleton';

    import { type application, type types, enums } from '$lib/wailsjs/go/models';
	import { UpdateProject } from '$lib/wailsjs/go/application/Driver';

    import { t } from '$lib/translations/translations';
    import { getSelectedProjectStore } from '$lib/stores';
    import { formatDateTime } from '$lib/components/utils';
    import { InputInline } from '$lib/components/ui/input';
    import { Gallery, Media, type GalleryItem } from '$lib/components/ui/gallery';

    import { AlertBuildRequired, AlertEmptyMedia, AlertMissingDependency } from './(components)/alert';

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

    function handleChange(event: CustomEvent) {
        updateProject();
    }

    function handleGalleryClick(event: CustomEvent<{ value: string }>) {
        if (!data.project.media) {
            return;
        }

        const filepath = event.detail.value;
        const index = data.project.media.findIndex((m) => m.filePath === filepath);

        if (index === -1) {
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

    $: thumbnail = data.project.media?.find((m) => m.thumbnail);
    $: build = data.project.dependencies?.find((d) => d.type === enums.PackageType.BUILD);

    $: dependencies = data.project.dependencies?.map((d) => {
        const pack = data.dependencies?.find((dep) => dep.reference === d.reference);

        return {
            id: pack?.id.toString() || "",
            reference: d.reference.toString(),
            installed: pack?.state === enums.PackageState.INSTALLED,
        }
    }) || [];

    $: dependenciesLabel = t.get('home.project.tag.dependency', { number: data.project.dependencies?.length || 0 });

    $: updatedAt = formatDateTime(data.project.updatedAt);
    $: galleryItems = convertToGalleryItems(data.project.media || []);
</script>

<main class="flex flex-col h-full space-y-4"> 
    <div class="flex gap-4 items-end">
        <div>
            <Media src={thumbnail?.url} height={32} width={32} class="cursor-default" rounded/>
        </div>
        <div class="space-y-2">
            <InputInline bind:value={data.project.name} labelClasses="h2 font-bold items-baseline" inputClasses="input" on:change={handleChange}>
                <IconEditFill class="text-sm text-surface-600-300-token"/>
            </InputInline>
            <div class="flex flex-wrap text-sm text-surface-800-100-token gap-1">
                <div class="badge variant-ghost rounded">{data.project.path}</div>
                <div class="badge variant-ghost rounded">{data.project.id}</div>
                <div class="badge variant-ghost rounded">{data.project.fileName}</div>
                {#if build}
                    <div class="badge variant-ghost rounded">{build.reference}</div>
                {/if}
                {#each data.project.tags || [] as tag}
                    <div class="badge variant-ghost-primary rounded">{tag}</div>
                {/each}
                {#if !build}
                    <div class="badge variant-ghost-error rounded">{dependenciesLabel}</div>
                {:else if dependencies?.find((d) => !d.installed)}
                    <div class="badge variant-ghost-warning rounded">{dependenciesLabel}</div>
                {:else}
                    <div class="badge variant-ghost-success rounded">{dependenciesLabel}</div>
                {/if}
                <div class="badge variant-ghost rounded">{updatedAt}</div>
            </div>
        </div>
    </div>
    <hr>
    <div class="h-full overflow-auto">
        <div class="px-2 space-y-4">
            {#if !build}
                <AlertBuildRequired />
            {/if}
            {#each dependencies as dependency}
                {#if !dependency.installed} 
                    <AlertMissingDependency
                        reference={dependency.reference}
                        id={dependency.id}
                    />
                {/if}
            {/each}
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
            <p class="text-xs font-semibold text-surface-600-300-token">Take care adding large files as it can cause performace issues.</p>
        </div>
    </div>
</main>