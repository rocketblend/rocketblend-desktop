<script lang="ts">
    import { onMount } from 'svelte';
    import SidebarHeader from '$lib/components/sidebar/SidebarHeader.svelte';
    import SidebarFilter from '$lib/components/sidebar/SidebarFilter.svelte';
    import PackageListBox from '$lib/components/package/PackageListBox.svelte';

    import type { packageservice } from '$lib/wailsjs/go/models';
    import { ListPackages } from '$lib/wailsjs/go/application/Driver';
    import { t } from '$lib/translations/translations';
    import type { RadioOption } from '$lib/types';
  
    let selectedFilterType: number = 0;
    let searchQuery: string = "";
    let filterInstalled: boolean = false;
    const filterRadioOptions: RadioOption[] = [
        { value: 0, key: 'all', display: $t('home.sidebar.filter.option.all') },
        { value: 1, key: 'build', display: $t('home.sidebar.filter.option.build') },
        { value: 2, key: 'addon', display: $t('home.sidebar.filter.option.addon') },
    ];

    let fetchPackagesPromise: Promise<packageservice.ListPackagesResponse | undefined>;
    let selectedPackageIds: string[] = [];

    async function fetchPackages(query: string): Promise<packageservice.ListPackagesResponse | undefined> {
        try {
            var category = filterRadioOptions[selectedFilterType].key;
            if (category === 'all') {
                category = '';
            }

            return await ListPackages(query, category, filterInstalled);
        } catch (error) {
            console.error('Error fetching packages:', error);
            return undefined;
        }
    }
  
    function handleInputChange(): void {
        fetchPackagesPromise = fetchPackages(searchQuery);
    }

    function handleAddPackage(): void {
        console.log('Add package');
    }

    function handleRefreshPackages(): void {
        console.log('Refresh packages');
    }
  
    onMount(() => {
        fetchPackagesPromise = fetchPackages(searchQuery);
    });
</script>
  
<div class="flex flex-col h-full space-y-4">
    <SidebarHeader 
        title={$t('home.sidebar.title')}
        on:add={handleAddPackage}
        on:refresh={handleRefreshPackages}
    />
    <SidebarFilter
        searchPlaceholder={$t('home.sidebar.filter.search')}
        installedLabel={$t('home.sidebar.filter.installed')}
        filterRadioOptions={filterRadioOptions}
        bind:selectedFilterType={selectedFilterType}
        bind:searchQuery={searchQuery}
        bind:filterInstalled={filterInstalled}
        on:filterChange={handleInputChange}
    />
    <div class="overflow-y-auto h-full">
        {#await fetchPackagesPromise}
            <div class="space-y-4 p-2">
                {#each Array(10) as _}
                    <div class="placeholder animate-pulse p-5 h-10" />
                {/each}
            </div>
        {:then response}
            {#if response && response.packages}
                <PackageListBox packages={response.packages} selectedPackageIds={selectedPackageIds} />
            {:else}
                <div class="p-2">
                    <p class="font-bold text-sm text-surface-200 text-center">{$t('home.sidebar.noresults')}</p>
                </div>
            {/if}
        {:catch error}
            <p>{$t('home.sidebar.error')}</p>
        {/await}
    </div>
</div>