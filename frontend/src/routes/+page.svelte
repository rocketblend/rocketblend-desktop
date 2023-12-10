<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { t } from '$lib/translations/translations';
    import { selectedProjectIds } from '$lib/store';
    import { RunProject } from '$lib/wailsjs/go/application/Driver';
    import type { project } from '$lib/wailsjs/go/models';
    import type { PageData } from './$types';
    import SearchFilter from '$lib/components/project/ProjectSearchFilter.svelte';
    import ProjectDisplay from '$lib/components/project/ProjectDisplay.svelte';

    export let data: PageData;

    let searchQuery = "";
    let displayType = "table";
    let form : HTMLFormElement;

    function handleSelected(event: CustomEvent<project.Project | null>): void {
        const project = event.detail;
        if (!project || !project.id) {
            return;
        }

        selectedProjectIds.set([project.id.toString()]);
    }

    function handleProjectDoubleClick(event: CustomEvent<{ itemId: string }>) {
        goto(`/projects/${event.detail.itemId}`);
    }

    function handleProjectActionDoubleClick(event: CustomEvent<{ itemId: string }>) {
        RunProject(event.detail.itemId)
    }

    function handleFilterChangeEvent(event: Event): void {
        form.requestSubmit();
    }

    $: searchQuery = $page.url.searchParams.get('query') || '';
    $: displayType = $page.url.searchParams.get('display') || 'table';
</script>

<main class="space-y-4"> 
    <h2 class="font-bold">{$t('home.title')}</h2>
    <div class="space-y-4">
        <SearchFilter
            bind:form={form}
            bind:searchQuery={searchQuery}
            bind:displayType={displayType}
            on:change={handleFilterChangeEvent} />
        <ProjectDisplay
            bind:projects={data.projects}
            bind:displayType={displayType}
            bind:selectedProjectIds={$selectedProjectIds}
            on:ctrlItemDoubleClicked={handleProjectActionDoubleClick}
            on:itemDoubleClicked={handleProjectDoubleClick}
            on:selected={handleSelected}/>
    </div>
</main>