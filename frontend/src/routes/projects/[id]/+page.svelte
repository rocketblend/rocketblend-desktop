<script lang="ts">
    import type { PageData } from './$types';
    import { FxReveal as Img } from '@zerodevx/svelte-img'

    import { selectedProject } from '$lib/store';

    import logo from '$lib/assets/images/logo.png?as=run';
	import Media from '$lib/components/media/Media.svelte';

    export let data: PageData;

    selectedProject.set(data.project);
</script>

<main class="space-y-4"> 
    <div class="flex gap-4">
        {#if data.project.thumbnailPath}
            <img class="h-auto max-w-full rounded-lg h-32 w-32" src="/system/{data.project.thumbnailPath}" alt="" loading="lazy" />
        {:else}
            <Img class="h-auto max-w-full rounded-lg h-32 w-32 image-fade-in" src={logo} alt="" loading="lazy" />
        {/if}
        <div class="relative w-full">
            <div class="absolute inset-x-0 bottom-0 space-y-4">
                <h2 class="font-bold">{data.project.name}</h2>
                <p class="text-sm text-surface-300">Last updated: {data.project.updatedAt}</p>
            </div>
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

    {#if data.project.thumbnailPath}
        <hr>
        <Media mediaClass="h-auto max-w-full rounded-lg" src="/system/{data.project.thumbnailPath}" alt="" />
    {/if}

    {#if data.project.splashPath}
        <hr>
        <Media mediaClass="h-auto max-w-full rounded-lg" src="/system/{data.project.splashPath}" alt="" />
    {/if}
</main>