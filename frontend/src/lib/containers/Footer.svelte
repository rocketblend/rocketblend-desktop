<script lang="ts">
    import { selectedProject } from '$lib/store';
    import { goto } from '$app/navigation';
    import { ExploreProject, RunProject } from '$lib/wailsjs/go/application/Driver';

    import FooterContent from '$lib/components/footer/FooterContent.svelte';
  
    function handleViewProject() {
        goto(`/projects/${$selectedProject?.id}`);
    }
  
    async function handleRunProject() {
        return await RunProject($selectedProject?.id);
    }
  
    async function handleExploreProject() {
        return await ExploreProject($selectedProject?.id);
    }
</script>
  
<FooterContent
    {selectedProject}
    onViewProject={handleViewProject}
    onRunProject={handleRunProject}
    onExploreProject={handleExploreProject}
/>