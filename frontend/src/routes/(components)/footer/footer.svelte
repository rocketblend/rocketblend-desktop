<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';

    import { EventsOn } from '$lib/wailsjs/runtime';
    import type { project } from '$lib/wailsjs/go/models';
    import { GetProject, ExploreProject, RunProject } from '$lib/wailsjs/go/application/Driver';

    import { resourcePath } from '$lib/components/utils';
    import { getSelectedProjectStore } from '$lib/stores';
    import { SEARCH_STORE_INSERT_CHANNEL } from '$lib/events';
    import { debounce } from '$lib/utils';

    import FooterContent from './footer-content.svelte';

    const selectedProjectStore = getSelectedProjectStore();
    const refreshProjectDebounced = debounce(loadProject, 1000);

    export let selectedProject: project.Project | undefined;

    let cancelListener: () => void;

    $: if ($selectedProjectStore) {
        loadProject();
    }

    async function loadProject() {
        var id = selectedProjectStore.latest();
        if (!id) {
            return;
        }

        const result = await GetProject(id);
        selectedProject = result.project;
    }

    function handleViewProject() {
        if (selectedProject) {
            goto(`/projects/${selectedProject.id}`);
        }
    }

    async function handleRunProject() {
        if (selectedProject) {
            return await RunProject(selectedProject.id);
        }
    }

    async function handleExploreProject() {
        if (selectedProject) {
            return await ExploreProject(selectedProject.id);
        }
    }

    onMount(() => {
        cancelListener = EventsOn(SEARCH_STORE_INSERT_CHANNEL, (event: { id: string, indexType: string }) => {
            if (selectedProject && event.indexType === "project" && selectedProject.id?.toString() === event.id) {
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

<FooterContent
    name={selectedProject?.name}
    fileName={selectedProject?.fileName}
    imagePath={resourcePath(selectedProject?.thumbnailPath)}
    isLoading={!selectedProject}
    on:viewProject={handleViewProject}
    on:runProject={handleRunProject}
    on:exploreProject={handleExploreProject}
/>