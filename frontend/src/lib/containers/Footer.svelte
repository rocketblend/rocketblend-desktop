<script lang="ts">
    import { selectedProjectIds } from '$lib/store';
    import { goto } from '$app/navigation';
    import { ExploreProject, RunProject } from '$lib/wailsjs/go/application/Driver';
    import { GetProject } from '$lib/wailsjs/go/application/Driver'
    import type { project } from '$lib/wailsjs/go/models';
    import { resourcePath } from '$lib/components/utils';
    import FooterContent from '$lib/components/footer/FooterContent.svelte';

    let selectedProject: project.Project | undefined = undefined;

    $: if ($selectedProjectIds) {
        loadProject();
    }

    async function loadProject() {
        var id = selectedProjectIds.latest();
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
    on:viewProject={handleViewProject}
    on:runProject={handleRunProject}
    on:exploreProject={handleExploreProject}
/>