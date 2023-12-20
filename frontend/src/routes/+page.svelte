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
	import { DisplayType, type OptionGroup } from '$lib/types';
	import { convertToEnum } from '$lib/components/utils';
    
    export let data: PageData;

    let displayTypeParam = "";
    let sortByParam = "";

    let searchQuery = "";
    let displayType: DisplayType = DisplayType.Table;

    let form : HTMLFormElement;

    const optionGroups: OptionGroup[] = [
        {
            label: 'sort',
            display: t.get('home.project.filter.group.sort.title'),
            options: [
                { value: 0, display: t.get('home.project.filter.group.sort.option.name') },
                { value: 1, display: t.get('home.project.filter.group.sort.option.file') },
                { value: 2, display: t.get('home.project.filter.group.sort.option.build') }
            ]
        }
    ];

    let primaryOptionGroup: number = 0;
    let selectedOptions: Record<string, number> = {'sort': 0};
    let optionLabel: string = t.get('home.project.filter.group.title');

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

    $: displayType = convertToEnum(displayTypeParam, DisplayType) || DisplayType.Table;

    $: optionLabel = optionGroups[primaryOptionGroup].options[selectedOptions[optionGroups[primaryOptionGroup].label]].display;
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
            bind:selectedOptions={selectedOptions}
            bind:filterLabel={optionLabel}
            optionsGroups={optionGroups}
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