<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { RadioGroup, RadioItem, SlideToggle } from '@skeletonlabs/skeleton';
    import type { RadioOption } from '$lib/types';

    import { InputSearch } from '$lib/components/ui/input';
	import type { enums } from '$lib/wailsjs/go/models';

    export let selectedFilterType: enums.PackageType | undefined;
    export let searchQuery: string;
    export let filterInstalled: boolean;
    export let filterRadioOptions: RadioOption[] = [];
    export let searchPlaceholder: string;
    export let installedLabel: string;

    const dispatch = createEventDispatcher();

    function handleFilterChange() {
        dispatch('filterChange');
    }
</script>

<div class="flex flex-col space-y-4">
    {#if filterRadioOptions.length > 0}
        <RadioGroup display="inline-flex">
            {#each filterRadioOptions as option}
                <RadioItem bind:group={selectedFilterType} name="justify" value={option.value} class="text-sm" on:change={handleFilterChange}>
                    {option.display}
                </RadioItem>
            {/each}
        </RadioGroup>
    {/if}
    <InputSearch bind:value={searchQuery} placeholder={searchPlaceholder} debounceDelay={500} on:input={handleFilterChange} class="text-sm"/>
    <SlideToggle name="slider-label" size="sm" active="bg-surface-200" class="text-sm" border="ring-outline-token" bind:checked={filterInstalled} on:change={handleFilterChange}>
        {installedLabel}
    </SlideToggle>
</div>