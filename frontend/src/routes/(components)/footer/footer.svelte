<script lang="ts">
    import { goto } from '$app/navigation';

    import type { project } from '$lib/wailsjs/go/models';
    import { ExploreProject, RunProject } from '$lib/wailsjs/go/application/Driver';

    import { resourcePath } from '$lib/components/utils';

    import FooterContent from './footer-content.svelte';

    export let selected: project.Project | undefined;

    function handleViewProject() {
        if (selected) {
            goto(`/projects/${selected.id}`);
        }
    }

    async function handleRunProject() {
        if (selected) {
            return await RunProject(selected.id);
        }
    }

    async function handleExploreProject() {
        if (selected) {
            return await ExploreProject(selected.id);
        }
    }
</script>

<FooterContent
    name={selected?.name}
    fileName={selected?.fileName}
    imagePath={resourcePath(selected?.thumbnailPath)}
    isLoading={!selected}
    on:viewProject={handleViewProject}
    on:runProject={handleRunProject}
    on:exploreProject={handleExploreProject}
/>