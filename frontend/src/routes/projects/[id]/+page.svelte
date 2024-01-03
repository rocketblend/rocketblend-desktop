<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import type { PageData } from './$types';

    import { getSelectedProjectStore } from '$lib/stores';
    import { debounce } from '$lib/utils';
    import { EventsOn } from '$lib/wailsjs/runtime';
    import { resourcePath } from '$lib/components/utils';

	import Media from '$lib/components/core/media/Media.svelte';
	import { GetProject } from '$lib/wailsjs/go/application/Driver';

    const selectedProjectStore = getSelectedProjectStore();
    const refreshProjectDebounced = debounce(refreshProject, 1000);

    export let data: PageData;
    
    let cancelListener: () => void;

    async function refreshProject() {
        const project = (await GetProject(data.project.id?.toString())).project;
        if (!project) {
            return;
        }

        data = {...data, project};
    }

    function setSelectedProject() {
        if (data.project.id) {
            selectedProjectStore.set([data.project.id.toString()]);
        }
    }

    setSelectedProject();

    onMount(() => {
        cancelListener = EventsOn('searchstore.insert', (event: { id: string, indexType: string }) => {
            if (event.indexType === "project" && event.id === data.project.id?.toString()) {
                refreshProjectDebounced();
            }
        });
    });

    onDestroy(() => {
        if (cancelListener) {
            cancelListener();
        }
    });
</script>

<main class="space-y-4"> 
    <div class="flex gap-4">
        <Media src={resourcePath(data.project.thumbnailPath)} alt="" />
        <div class="space-y-4">
            <h2 class="h2 font-bold">{data.project.name}</h2>
            <p class="text-sm text-surface-300">Last updated: {data.project.updatedAt}</p>
        </div>
    </div>
    <hr>
    <ul>
        <li>ID: {data.project.id}</li>
        <li>Path: {data.project.path}</li>
        <li>File Name: {data.project.fileName}</li>
        <li>Thumbnail Path: {data.project.thumbnailPath}</li>
        <li>Splash Path: {data.project.splashPath}</li>
        <li>Build: {data.project.build}</li>
        <li>Addons: {data.project.addons}</li>
        <li>Tags: {data.project.tags}</li>
        <li>Version: {data.project.version}</li>
    </ul>
    <hr>
    <div class="grid grid-cols-4 gap-4">
        <Media height="80" width="full" src="{resourcePath(data.project.splashPath)}" alt="" />
      </div>
</main>