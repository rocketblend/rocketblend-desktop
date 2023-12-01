<script lang="ts">
    import { onMount } from 'svelte';

    import { RadioGroup, RadioItem, SlideToggle } from '@skeletonlabs/skeleton';

    import IconBox2Fill from '~icons/ri/box-2-fill'

    import { t } from '$lib/translations/translations';
    import type { packageservice } from '$lib/wailsjs/go/models';
    import { ListPackages } from '$lib/wailsjs/go/application/Driver'
    import SearchInput from '$lib/components/core/search-input/search-input.svelte';

    let filterType: number;
    let query: string = "";
    let fetchPackagesPromise : Promise<packageservice.ListPackagesResponse| undefined> ;

    type RadioOption = {
        value: number;
        label: string;
    };

    const radioOptions: RadioOption[] = [
        { value: 0, label: 'All' },
        { value: 1, label: 'Builds' },
        { value: 2, label: 'Addons' },
    ];

    async function fetchPackages(query:string): Promise<packageservice.ListPackagesResponse | undefined> {
        try {
            return await ListPackages(query);
        } catch (error) {
            console.error('Error fetching packages:', error);
            return undefined;
        }
    }

    function handleInputChange(event: Event): void {
        fetchPackagesPromise = fetchPackages(query);
    }

    onMount(() => {
        fetchPackagesPromise = fetchPackages(query);
    });
</script>

<div class="flex flex-col h-full space-y-4">
    <div class="inline-flex items-center space-x-2 text-surface-200">
        <IconBox2Fill/>
        <h5 class="font-bold">{$t('home.sidebar.title')}</h5>
    </div>
    <RadioGroup display="inline-flex">
        {#each radioOptions as option}
          <RadioItem bind:group={filterType} name="justify" value={option.value} class="text-sm">{option.label}</RadioItem>
        {/each}
    </RadioGroup>
    <SearchInput bind:value={query} placeholder={$t('home.sidebar.search')} debounceDelay={500} on:input={handleInputChange} class="text-sm"/>
    <SlideToggle name="slider-label" size="sm" active="bg-surface-200" class="text-sm">{$t('home.sidebar.installed')}</SlideToggle>
    <div class="overflow-y-auto h-full">
        {#await fetchPackagesPromise}
            <div class="flex-auto space-y-4 p-2">
                {#each Array(10) as _}
                    <div class="placeholder animate-pulse p-5" />
                {/each}
            </div>
        {:then response}
            {#if response && response.packages}
                <dl class="flex-auto list-dl">
                    {#each response.packages || [] as pack }
                        <div>
                            <span class="flex-auto text-ellipsis overflow-hidden">
                                <dt class="font-bold text-sm">{pack.name}</dt>
                                <dd class="text-surface-300 text-xs">{pack.reference}</dd>
                            </span>
                        </div>
                    {/each}
                </dl>
            {:else}
                <div class="flex-auto p-2">
                    <p class="font-bold text-sm text-surface-200 text-center">No packages found!</p>
                </div>
            {/if}
        {:catch error}
            <p>An error occurred while fetching packages!</p>
        {/await}
    </div>
</div>
