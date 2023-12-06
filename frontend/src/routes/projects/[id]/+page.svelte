<script lang="ts">
    import type { PageData } from './$types';

    import { resourcePath } from '$lib/components/utils';
    import { selectedProject } from '$lib/store';

	import Media from '$lib/components/core/media/Media.svelte';

    export let data: PageData;

    selectedProject.set(data.project);
</script>

<main class="space-y-4"> 
    <div class="flex gap-4">
        <Media src={resourcePath(data.project.thumbnailPath)} alt="" />
        <div class="space-y-4">
            <h2 class="font-bold">{data.project.name}</h2>
            <p class="text-sm text-surface-300">Last updated: {data.project.updatedAt}</p>
        </div>
    </div>
    <hr>
    <ul>
        <li>ID: {data.project.id}</li>
        <li>Path: {data.project.path}</li>
        <li>File Name: {data.project.fileName}</li>
        <li>Thumbnail Path: {data.project.thumbnailPath}</li>
        <li>Splash Path: {data.project.splashPath}</li>
        <li>Build: {data.project.build}</li>
        <li>Addons: {data.project.addons}</li>
        <li>Tags: {data.project.tags}</li>
        <li>Version: {data.project.version}</li>
    </ul>
    <hr>
    <div class="grid grid-cols-4 gap-4">
        <Media height="auto" width="full" src="{resourcePath(data.project.splashPath)}" alt="" />
      </div>
</main>