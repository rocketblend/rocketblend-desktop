<script lang="ts">
    import type { PageData } from './$types';
    import { goto, invalidate } from '$app/navigation';
    import { page } from '$app/stores';

    import { t } from '$lib/translations/translations';
    import { RunProject } from '$lib/wailsjs/go/application/Driver';

    import { getSelectedProjectStore } from '$lib/stores';
	import { DisplayType, type OptionGroup } from '$lib/types';
	import { convertToEnum, debounce } from '$lib/components/utils';

    import {
        ProjectList,
        ProjectFilter,
        ProjectCreateButton,
        AlertNoProjectDirectory
    } from "./(components)"

    const selectedProjectStore = getSelectedProjectStore();

    // const optionGroups: OptionGroup[] = [
    //     {
    //         label: 'sort',
    //         display: t.get('home.project.filter.group.sort.title'),
    //         options: [
    //             { value: 0, display: t.get('home.project.filter.group.sort.option.name') },
    //             { value: 1, display: t.get('home.project.filter.group.sort.option.file') },
    //             { value: 2, display: t.get('home.project.filter.group.sort.option.build') }
    //         ]
    //     }
    // ];
    
    export let data: PageData;

    let enabled = data.preferences.watchPath !== "";

    let displayTypeParam = "";
    let sortByParam = "";

    let searchQuery = "";
    let displayType: DisplayType = DisplayType.Table;

    let form : HTMLFormElement;

    // let primaryOptionGroup: number = 0;
    // let selectedOptions: Record<string, number> = {'sort': 0};
    // let optionLabel: string = t.get('home.project.filter.group.title');

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

    $: enabled = data.preferences.watchPath !== "";
    $: searchQuery = $page.url.searchParams.get("query") || "";
    $: displayTypeParam = $page.url.searchParams.get("display") || "";
    $: sortByParam = $page.url.searchParams.get("sortBy") || "";

    $: displayType = convertToEnum(displayTypeParam, DisplayType) || DisplayType.Table;

    $: $selectedProjectStore ? invalidate("app:layout"): null;

    // $: optionLabel = optionGroups[primaryOptionGroup].options[selectedOptions[optionGroups[primaryOptionGroup].label]].display;
</script>

<main class="flex flex-col h-full space-y-4">
    <div class="flex justify-between items-center">
        <h2 class="h2 font-bold">{$t('home.title')}</h2>
        <ProjectCreateButton disabled={!enabled}/>
    </div>
    <hr>
    {#if enabled}
        <ProjectFilter
            bind:form={form}
            bind:searchQuery={searchQuery}
            bind:displayType={displayType}
            searchPlaceholder={$t('home.project.query.placeholder')}
            on:change={handleFilterChangeEvent} />
        <ProjectList
            bind:projects={data.projects}
            bind:displayType={displayType}
            bind:selectedProjectIds={$selectedProjectStore}
            on:itemDoubleClick={handleProjectDoubleClick}
            on:sortChanged={handleSortChange}/>
    {:else}
        <AlertNoProjectDirectory />
    {/if}
</main>