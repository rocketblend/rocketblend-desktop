<script lang="ts">
    import { goto } from '$app/navigation';

    import { getDrawerStore } from '@skeletonlabs/skeleton';

    import { GetProject, ExploreProject, RunProject } from '$lib/wailsjs/go/application/Driver';
    import type { project } from '$lib/wailsjs/go/models';

    import { resourcePath } from '$lib/components/utils';
    import { getSelectedProjectStore } from '$lib/stores';

    import FooterContent from '$lib/components/footer/FooterContent.svelte';

    const selectedProjectStore = getSelectedProjectStore();
    const drawerStore = getDrawerStore();

    let selectedProject: project.Project | undefined = undefined;

    $: if ($selectedProjectStore) {
        loadProject();
    }

    function handleViewTerminal() {
        drawerStore.open();
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
</script>

<FooterContent
    name={selectedProject?.name}
    fileName={selectedProject?.fileName}
    imagePath={resourcePath(selectedProject?.thumbnailPath)}
    isLoading={!selectedProject}
    on:viewTerminal={handleViewTerminal}
    on:viewProject={handleViewProject}
    on:runProject={handleRunProject}
    on:exploreProject={handleExploreProject}
/>