<script lang="ts">
    import { onMount } from 'svelte';
    
	import ProjectGallery from '../project-gallery/project-gallery.svelte';
    import ProjectTable from '../project-table/project-table.svelte';
    import SearchInput from '../core/search-input/search-input.svelte';

    import type { project } from '$lib/wailsjs/go/models';
    import { FindAllProjects } from '$lib/wailsjs/go/application/Driver';

    let projects: project.Project[] = [];
    let searchTerm: string = '';

    async function loadProjects(query: string): Promise<void> {
        const result = await FindAllProjects(query);
        projects = result || [];
    }

    function handleSearch(event: CustomEvent<string>): void {
        searchTerm = event.detail;
    }

    // Load initial project list
    onMount(async () => {
        await loadProjects(searchTerm);
    });

    // Reactive statement to fetch new project list when searchTerm changes
    $: loadProjects(searchTerm), searchTerm;
</script>

<div class="space-y-4">
    <div>
        <SearchInput on:search={handleSearch} placeholder="Search"/>
    </div>
    
    <ProjectTable sourceData={projects}/>
    <!-- <ProjectGallery/> -->
</div>
