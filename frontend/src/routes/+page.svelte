<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import type { PageData } from './$types';

    import { goto } from '$app/navigation';
    import { page } from '$app/stores';

    import { getToastStore } from '@skeletonlabs/skeleton';

    import { EventsOn } from '$lib/wailsjs/runtime';
    import { t } from '$lib/translations/translations';
    import { RunProject, ListProjects } from '$lib/wailsjs/go/application/Driver';

    import { getSelectedProjectStore } from '$lib/stores';
	import { DisplayType, type OptionGroup } from '$lib/types';
    import { changeDetectedToast } from '$lib/toasts';
	import { convertToEnum, debounce } from '$lib/components/utils';
    

    import ProjectListView from '$lib/components/project/ProjectListView.svelte';
	import ProjectFilter from '$lib/components/project/ProjectFilter.svelte';

    const selectedProjectStore = getSelectedProjectStore();
    const toastStore = getToastStore();
    const fetchProjectsDebounced = debounce(refreshProjects, 500);
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
    
    export let data: PageData;

    let displayTypeParam = "";
    let sortByParam = "";

    let searchQuery = "";
    let displayType: DisplayType = DisplayType.Table;

    let form : HTMLFormElement;

    let primaryOptionGroup: number = 0;
    let selectedOptions: Record<string, number> = {'sort': 0};
    let optionLabel: string = t.get('home.project.filter.group.title');

    let cancelListener: () => void;

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

    async function refreshProjects() {
        toastStore.trigger(changeDetectedToast);
        const projects = (await ListProjects(searchQuery)).projects;
        data = {...data, projects};
    }

    onMount(() => {
        cancelListener = EventsOn('searchstore.insert', (data: { id: string, indexType: string }) => {
            if (data.indexType === "project") {
                fetchProjectsDebounced();
            }
        });
    });

    onDestroy(() => {
        if (cancelListener) {
            cancelListener();
        }
    });

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
            bind:selectedProjectIds={$selectedProjectStore}
            on:itemDoubleClick={handleProjectDoubleClick}
            on:sortChanged={handleSortChange}/>
    </div>
</main>