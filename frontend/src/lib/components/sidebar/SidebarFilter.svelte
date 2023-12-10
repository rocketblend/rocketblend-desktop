<script lang="ts">
    import { RadioGroup, RadioItem, SlideToggle } from '@skeletonlabs/skeleton';

    import { t } from '$lib/translations/translations';

    import type { RadioOption } from '$lib/types';

    import SearchInput from '$lib/components/core/input/SearchInput.svelte';

    export let selectedFilterType: number;
    export let searchQuery: string;
    export let filterInstalled: boolean;
    export let filterRadioOptions: RadioOption[];
    export let onFilterChange: () => void;
</script>
  
<div class="flex flex-col space-y-4">
    <RadioGroup display="inline-flex">
        {#each filterRadioOptions as option}
            <RadioItem bind:group={selectedFilterType} name="justify" value={option.value} class="text-sm" on:change={onFilterChange}>
                {$t(`home.sidebar.filter.option.${option.key}`)}
            </RadioItem>
        {/each}
    </RadioGroup>
    <SearchInput bind:value={searchQuery} placeholder={$t('home.sidebar.filter.search')} debounceDelay={500} on:input={onFilterChange} class="text-sm"/>
    <SlideToggle name="slider-label" size="sm" active="bg-surface-200" class="text-sm" border="ring-outline-token" bind:checked={filterInstalled} on:change={onFilterChange}>
        {$t('home.sidebar.filter.installed')}
    </SlideToggle>
</div>