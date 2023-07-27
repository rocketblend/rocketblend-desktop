<script lang="ts">
	import ProjectGallery from '../project-gallery/project-gallery.svelte';
    import ProjectTable from '../project-table/project-table.svelte';

    import type { project } from '$lib/wailsjs/go/models';

    export let projects: project.Project[];
</script>


<ProjectTable sourceData={projects}/>
<!-- <ProjectGallery/> -->