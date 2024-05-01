<script lang="ts">
    import { goto } from '$app/navigation';

    import { type types, application } from '$lib/wailsjs/go/models';
    import { RunProject, OpenExplorer } from '$lib/wailsjs/go/application/Driver';

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
            const opts = application.OpenExplorerOptions.createFrom({
                path: selected.path,
            });
            
            await OpenExplorer(opts);
        }
    }
</script>

<FooterContent
    name={selected?.name}
    fileName={selected?.fileName}
    imagePath={selected?.thumbnail?.url}
    isLoading={!selected}
    on:viewProject={handleViewProject}
    on:runProject={handleRunProject}
    on:exploreProject={handleExploreProject}
/>