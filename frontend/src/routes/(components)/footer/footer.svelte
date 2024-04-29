<script lang="ts">
    import { goto } from '$app/navigation';

    import { type types, application } from '$lib/wailsjs/go/models';
    import { RunProject } from '$lib/wailsjs/go/application/Driver';

    import { resourcePath } from '$lib/components/utils';

    import FooterContent from './footer-content.svelte';

    export let selected: types.Project | undefined;

    function handleViewProject() {
        if (selected) {
            goto(`/projects/${selected.id}`);
        }
    }

    async function handleRunProject() {
        if (selected) {
            const opts = application.RunProjectOpts.createFrom({
                id: selected.id,
            });
            
            return await RunProject(opts);
        }
    }

    async function handleExploreProject() {
        if (selected) {
            console.log("Explore project");
            // TODO: Switch to general explore endpoint.
            // return await ExploreProject(selected.id);
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