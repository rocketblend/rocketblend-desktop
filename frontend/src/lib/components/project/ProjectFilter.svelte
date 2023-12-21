<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import SearchInput from '$lib/components/core/input/SearchInput.svelte';
    import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';

    import IconListUnordered from '~icons/ri/list-unordered';
    import IconGalleryView from '~icons/ri/gallery-view-2';
    import IconTableView from '~icons/ri/table-view';

    import { DisplayType, type OptionGroup } from '$lib/types';
	import OptionButton from '../button/OptionButton.svelte';

    export let form: HTMLFormElement;
    export let searchQuery: string;
    export let searchPlaceholder: string = "Search";
    export let displayType: DisplayType = DisplayType.Table;

    export let optionsGroups: OptionGroup[] = [];
    export let selectedOptions: Record<string, number> = {};
    export let filterLabel: string = "Filter";

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change')
    }
</script>

<form bind:this={form} class="inline-flex space-x-2 w-full" data-sveltekit-keepfocus>
    <div class="flex-grow">
        <SearchInput name="query" value={searchQuery} on:input={handleChange} placeholder={searchPlaceholder} debounceDelay={500} />
    </div>
    <OptionButton {optionsGroups} bind:selectedOptions on:optionChange={handleChange}>
        <svelte:fragment slot="buttonContent">
            <IconListUnordered />
            <span class="capitalize">{filterLabel}</span>
        </svelte:fragment>
    </OptionButton>
    <RadioGroup display="flex items-center">
        <RadioItem bind:group={displayType} name="display" value={DisplayType.Table} on:change={handleChange}><IconTableView/></RadioItem>
        <RadioItem bind:group={displayType} name="display" value={DisplayType.Gallery} on:change={handleChange}><IconGalleryView/></RadioItem>
    </RadioGroup>
    <button type="submit" class="hidden"></button>
</form>