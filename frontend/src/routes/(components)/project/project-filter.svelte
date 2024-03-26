<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';

    // import IconListUnordered from '~icons/ri/list-unordered';
    import IconGalleryView from '~icons/ri/gallery-view-2';
    import IconTableView from '~icons/ri/table-view';

    import { DisplayType, type OptionGroup } from '$lib/types';

	// import { ButtonOption } from "$lib/components/ui/button"
    import { InputSearch } from "$lib/components/ui/input"

    export let form: HTMLFormElement;
    export let searchQuery: string;
    export let searchPlaceholder: string = "Search";
    export let displayType: DisplayType = DisplayType.Table;

    // export let optionsGroups: OptionGroup[] = [];
    // export let selectedOptions: Record<string, number> = {};
    // export let filterLabel: string = "Filter";

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change')
    }
</script>

<form bind:this={form} class="inline-flex space-x-2 w-full" data-sveltekit-keepfocus>
    <div class="flex-grow">
        <InputSearch
            name="query"
            value={searchQuery}
            placeholder={searchPlaceholder}
            debounceDelay={500}
            on:input={handleChange}
        />
    </div>
    <!-- <ButtonOption {optionsGroups} bind:selectedOptions on:optionChange={handleChange}>
        <svelte:fragment slot="buttonContent">
            <IconListUnordered />
            <span class="capitalize">{filterLabel}</span>
        </svelte:fragment>
    </ButtonOption> -->
    <RadioGroup display="flex items-center">
        <RadioItem bind:group={displayType} name="display" value={DisplayType.Table} on:change={handleChange}><IconTableView/></RadioItem>
        <RadioItem bind:group={displayType} name="display" value={DisplayType.Gallery} on:change={handleChange}><IconGalleryView/></RadioItem>
    </RadioGroup>
    <button type="submit" class="hidden"></button>
</form>