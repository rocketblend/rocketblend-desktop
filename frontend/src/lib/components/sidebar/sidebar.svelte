<script lang="ts">
    import { onMount } from 'svelte';

    import { RadioGroup, RadioItem, SlideToggle, ListBox, ListBoxItem } from '@skeletonlabs/skeleton';

    import IconBox2Fill from '~icons/ri/box-2-fill'
    import IconLoopRightFill from '~icons/ri/loop-right-fill'
    import IconAddBoxFill from '~icons/ri/add-box-fill'
    import IconCheckboxBlankCircleLine from '~icons/ri/checkbox-blank-circle-line'
    import IconCheckboxBlankLine from '~icons/ri/checkbox-blank-line'

    import { t } from '$lib/translations/translations';
    import type { pack, packageservice } from '$lib/wailsjs/go/models';
    import { ListPackages } from '$lib/wailsjs/go/application/Driver'
    import SearchInput from '$lib/components/core/search-input/search-input.svelte';

    let filterType: number = 0;
    let query: string = "";
    let filerInstalled: boolean = false;
    let fetchPackagesPromise : Promise<packageservice.ListPackagesResponse| undefined> ;

    let selectedPackageIds: string[] = [];

    type RadioOption = {
        value: number;
        key: string;
    };

    const radioOptions: RadioOption[] = [
        { value: 0, key: 'all' },
        { value: 1, key: 'build' },
        { value: 2, key: 'addon' },
    ];

    async function fetchPackages(query:string): Promise<packageservice.ListPackagesResponse | undefined> {
        try {
            var category = radioOptions[filterType].key;
            if (category == 'all') {
                category = '';
            }

            return await ListPackages(query, category, filerInstalled);
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
    <div class="inline-flex items-center align-center space-x-2 text-surface-200">
        <IconBox2Fill/>
        <h5 class="font-bold flex-grow">{$t('home.sidebar.title')}</h5>
        <button type="button" class="btn p-0 text-surface-200" >
            <IconLoopRightFill class="text-md mt-1"/>
        </button>
        <button type="button" class="btn p-0 text-surface-200" >
            <IconAddBoxFill class="text-xl mt-1"/>
        </button>
    </div>
    <!-- <div>
        <ProgressBar meter="bg-primary-400-500-token" height="h-1"/>
    </div> -->
    <RadioGroup display="inline-flex">
        {#each radioOptions as option}
        <RadioItem bind:group={filterType} name="justify" value={option.value} class="text-sm" on:change={handleInputChange}>
            {$t(`home.sidebar.filter.option.${option.key}`)}
          </RadioItem>
        {/each}
    </RadioGroup>
    <SearchInput bind:value={query} placeholder={$t('home.sidebar.filter.search')} debounceDelay={500} on:input={handleInputChange} class="text-sm"/>
    <SlideToggle name="slider-label" size="sm" active="bg-surface-200" class="text-sm" border="ring-outline-token" bind:checked={filerInstalled} on:change={handleInputChange}>{$t('home.sidebar.filter.installed')}</SlideToggle>
    <div class="overflow-y-auto h-full">
        {#await fetchPackagesPromise}
            <div class="flex-auto space-y-4 p-2">
                {#each Array(10) as _}
                    <div class="placeholder animate-pulse p-5 h-10" />
                {/each}
            </div>
        {:then response}
            {#if response && response.packages}
            <ListBox class="flex-auto" multiple>
                {#each response.packages || [] as pack }
                    <ListBoxItem bind:group={selectedPackageIds} name="packages" value={pack.id} active="variant-glass-primary" hover="hover:variant-filled-surface" rounded="rounded" class="truncate overflow-hidden">
                        <div class="flex gap-2">
                            <!-- <div class="flex-shrink-0 inline-flex items-center">
                                {#if pack.type == 'build'}
                                    <IconCheckboxBlankCircleLine />
                                {:else}
                                    <IconCheckboxBlankLine />
                                {/if}
                            </div> -->
                            <div class="flex-grow w-2">
                                <div class="text-sm font-medium">{pack.name}</div>
                                <div class="text-sm">{pack.reference}</div>
                            </div>
                        </div>
                    </ListBoxItem>
                {/each}
            </ListBox>
            {:else}
                <div class="flex-auto p-2">
                    <p class="font-bold text-sm text-surface-200 text-center">{$t('home.sidebar.noresults')}</p>
                </div>
            {/if}
        {:catch error}
            <p>An error occurred while fetching packages!</p>
        {/await}
    </div>
</div>
