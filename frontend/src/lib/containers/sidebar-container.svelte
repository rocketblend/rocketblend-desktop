<script lang="ts">
    import { onMount } from 'svelte';
    import SidebarHeader from '$lib/components/sidebar-header/sidebar-header.svelte';
    import SidebarFilters from '$lib/components/sidebar-filters/sidebar-filters.svelte';
    import PackageList from '$lib/components/package-list/package-list.svelte';
    import type { pack, packageservice } from '$lib/wailsjs/go/models';
    import { ListPackages } from '$lib/wailsjs/go/application/Driver';
    import { t } from '$lib/translations/translations';
    import type { RadioOption } from '$lib/types/radio-option';
  
    let selectedFilterType: number = 0;
    let searchQuery: string = "";
    let filterInstalled: boolean = false;
    const filterRadioOptions: RadioOption[] = [
        { value: 0, key: 'all' },
        { value: 1, key: 'build' },
        { value: 2, key: 'addon' },
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
        onAdd={handleAddPackage}
        onRefresh={handleRefreshPackages}
    />
    <SidebarFilters
        bind:selectedFilterType={selectedFilterType}
        bind:searchQuery={searchQuery}
        bind:filterInstalled={filterInstalled}
        filterRadioOptions={filterRadioOptions}
        onFilterChange={handleInputChange}
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
                <PackageList packages={response.packages} bind:selectedPackageIds={selectedPackageIds} />
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