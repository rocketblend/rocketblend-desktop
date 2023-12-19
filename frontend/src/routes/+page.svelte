<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { t } from '$lib/translations/translations';
    import { selectedProjectIds } from '$lib/store';
    import { RunProject } from '$lib/wailsjs/go/application/Driver';
    import type { project } from '$lib/wailsjs/go/models';
    import type { PageData } from './$types';

	import ProjectListView from '$lib/components/project/ProjectListView.svelte';
	import ProjectFilter from '$lib/components/project/ProjectFilter.svelte';
	import { DisplayType, SortBy, type RadioOption } from '$lib/types';
	import { convertToEnum } from '$lib/components/utils';
    
    export let data: PageData;

    let displayTypeParam = "";
    let sortByParam = "";

    let searchQuery = "";
    let displayType: DisplayType = DisplayType.Table;

    let sortBy: number = SortBy.Name;
    let form : HTMLFormElement;

    const sortByOptions: RadioOption[] = [
        { value: SortBy.Name, display: t.get('home.project.filter.option.name') },
        { value: SortBy.File, display: t.get('home.project.filter.option.file') },
        { value: SortBy.Build, display: t.get('home.project.filter.option.build') },
    ];

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

    $: searchQuery = $page.url.searchParams.get("query") || "";
    $: displayTypeParam = $page.url.searchParams.get("display") || "";
    $: sortByParam = $page.url.searchParams.get("sortBy") || "";

    $: displayType = convertToEnum(displayTypeParam, DisplayType);
    //$: sortBy = convertToEnum(sortByParam, SortBy);
</script>

<main class="space-y-4">
    <div>
        <h2 class="h2 font-bold">{$t('home.title')}</h2>
    </div>

    <div class="space-y-4">
        <ProjectFilter
            bind:form={form}
            bind:searchQuery={searchQuery}
            bind:displayType={displayType}
            bind:sortBy={sortBy}
            sortByOptions={sortByOptions}
            searchPlaceholder={$t('home.project.query.placeholder')}
            on:change={handleFilterChangeEvent} />
        <ProjectListView
            bind:projects={data.projects}
            bind:displayType={displayType}
            bind:selectedProjectIds={$selectedProjectIds}
            on:ctrlItemDoubleClicked={handleProjectActionDoubleClick}
            on:itemDoubleClicked={handleProjectDoubleClick}
            on:selected={handleSelected}/>
    </div>
</main>