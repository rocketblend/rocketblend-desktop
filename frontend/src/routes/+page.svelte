<script lang="ts">
    import type { PageData } from './$types';

    import { goto } from '$app/navigation';
    import { page } from '$app/stores';

    import { t } from '$lib/translations/translations';
    import { selectedProjectIds } from '$lib/stores';
    import { RunProject } from '$lib/wailsjs/go/application/Driver';
	import { DisplayType, type OptionGroup } from '$lib/types';
	import { convertToEnum } from '$lib/components/utils';

    import ProjectListView from '$lib/components/project/ProjectListView.svelte';
	import ProjectFilter from '$lib/components/project/ProjectFilter.svelte';
	import LongRunningOperation from '$lib/containers/LongRunningOperation.svelte';
    
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

    function handleProjectDoubleClick(event: CustomEvent<{ event: MouseEvent, item: string }>) {
        if (event.detail.event.ctrlKey) {
            RunProject(event.detail.item)
            return;
        }

        goto(`/projects/${event.detail.item}`);
    }

    function handleFilterChangeEvent(event: Event): void {
        form.requestSubmit();
    }

    function handleSortChange(event: CustomEvent<{ key: string, direction: string }>) {
        return;
    }

    $: searchQuery = $page.url.searchParams.get("query") || "";
    $: displayTypeParam = $page.url.searchParams.get("display") || "";
    $: sortByParam = $page.url.searchParams.get("sortBy") || "";

    $: displayType = convertToEnum(displayTypeParam, DisplayType) || DisplayType.Table;

    $: optionLabel = optionGroups[primaryOptionGroup].options[selectedOptions[optionGroups[primaryOptionGroup].label]].display;
</script>

<main class="space-y-4">
    <LongRunningOperation />

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
            on:itemDoubleClick={handleProjectDoubleClick}
            on:sortChanged={handleSortChange}/>
    </div>
</main>