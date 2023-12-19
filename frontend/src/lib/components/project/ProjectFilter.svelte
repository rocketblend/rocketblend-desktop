<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import SearchInput from '$lib/components/core/input/SearchInput.svelte';
    import { popup } from '@skeletonlabs/skeleton';
    import { RadioGroup, RadioItem, ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
    import type { PopupSettings } from '@skeletonlabs/skeleton';

    import IconListUnordered from '~icons/ri/list-unordered';
    import IconGalleryView from '~icons/ri/gallery-view-2';
    import IconTableView from '~icons/ri/table-view';

    import { type RadioOption, DisplayType, SortBy } from '$lib/types';

    export let form: HTMLFormElement;
    export let searchQuery: string;
    export let searchPlaceholder: string = "Search";
    export let displayType: DisplayType = DisplayType.Table;
    export let sortBy: number = 0;
    export let sortByOptions: RadioOption[] = [];

    const dispatch = createEventDispatcher();

    const filterOptionsCombobox: PopupSettings = {
        event: 'focus-click',
        target: 'filterOptionsCombobox',
        placement: 'bottom',
        closeQuery: '.listbox-item'
    };

    function handleChange() {
        dispatch('change')
    }
</script>

<form bind:this={form} class="inline-flex space-x-2 w-full" data-sveltekit-keepfocus>
    <div class="flex-grow">
        <SearchInput name="query" value={searchQuery} on:input={handleChange} placeholder={searchPlaceholder} debounceDelay={500} />
    </div>
    <button class="btn justify-between variant-ghost-surface" use:popup={filterOptionsCombobox}>
        <span class="capitalize">{sortByOptions[sortBy].display ?? 'sort by'}</span>
        <IconListUnordered />
    </button>
    <div class="card w-48 p-2 variant-filled-surface rounded z-50" data-popup={filterOptionsCombobox.target}>
        <ListBox padding="px-2 py-2">
            <div class="pb-2 px-2 pt-1">
                <span class="text-xs font-bold text-surface-700-200-token">Sort By</span>
            </div>
            {#each sortByOptions as option}
                <ListBoxItem bind:group={sortBy} name="sort-by" value={option.value} on:change={handleChange}>{option.display}</ListBoxItem>
            {/each}
        </ListBox>
        <div class="arrow variant-filled-surface" />
    </div>
    <RadioGroup display="flex items-center">
        <RadioItem bind:group={displayType} name="display" value={DisplayType.Table} on:change={handleChange}><IconTableView/></RadioItem>
        <RadioItem bind:group={displayType} name="display" value={DisplayType.Gallery} on:change={handleChange}><IconGalleryView/></RadioItem>
    </RadioGroup>
    <button type="submit" class="hidden"></button>
</form>