<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import type { PageData } from './$types';

    import { getSelectedProjectStore } from '$lib/stores';
    import { debounce } from '$lib/utils';
    import { EventsOn } from '$lib/wailsjs/runtime';
    import { resourcePath } from '$lib/components/utils';
    import { EVENT_DEBOUNCE, SEARCH_STORE_INSERT_CHANNEL } from '$lib/events';

	import Media from '$lib/components/core/media/Media.svelte';
	import { GetProject } from '$lib/wailsjs/go/application/Driver';
	import InlineInput from '$lib/components/core/input/InlineInput.svelte';

    const selectedProjectStore = getSelectedProjectStore();
    const refreshProjectDebounced = debounce(refreshProject, EVENT_DEBOUNCE);

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
        cancelListener = EventsOn(SEARCH_STORE_INSERT_CHANNEL, (event: { id: string, indexType: string }) => {
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
    <div class="flex gap-4 items-end">
        <Media src={resourcePath(data.project.thumbnailPath)} alt="" />
        <div class="space-y-4">
            <InlineInput bind:value={data.project.name} labelClasses="h2 font-bold items-baseline" inputClasses="input" />
            <span class="text-sm text-surface-600-300-token">Last updated: {data.project.updatedAt}</span>
        </div>
    </div>
    <hr>
    <InlineInput type="textarea" placeholder="Description..."/>
    <hr>
    <ul>
        <li><b>ID:</b> {data.project.id}</li>
        <li><b>Path:</b> {data.project.path}</li>
        <li><b>File Name:</b> {data.project.fileName}</li>
        <li><b>Build:</b> {data.project.build}</li>
        <li><b>Addons:</b> {data.project.addons}</li>
        <li><b>Tags:</b> {data.project.tags}</li>
        <li><b>Version:</b> {data.project.version}</li>
    </ul>
    <hr>
    <div class="grid grid-cols-4 gap-4">
        <Media height="80" width="full" src="{resourcePath(data.project.splashPath)}" alt="" />
      </div>
</main>