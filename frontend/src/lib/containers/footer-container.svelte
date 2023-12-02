<script lang="ts">
    import Footer from '$lib/components/footer/footer.svelte';
    import { selectedProject } from '$lib/store';
    import { goto } from '$app/navigation';
    import { ExploreProject, RunProject } from '$lib/wailsjs/go/application/Driver';
  
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
  
<Footer
    {selectedProject}
    onViewProject={handleViewProject}
    onRunProject={handleRunProject}
    onExploreProject={handleExploreProject}
/>