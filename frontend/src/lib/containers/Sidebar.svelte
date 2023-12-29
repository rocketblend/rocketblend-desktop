<script lang="ts">
    import { onMount } from 'svelte';
    
    import type {  ToastSettings } from '@skeletonlabs/skeleton';
    import { getToastStore } from '@skeletonlabs/skeleton';

    import { t } from '$lib/translations/translations';

    import type { packageservice } from '$lib/wailsjs/go/models';
    import { GetProject, ListPackages } from '$lib/wailsjs/go/application/Driver';

    import type { RadioOption } from '$lib/types';
    import { getSelectedProjectStore } from '$lib/stores';

    import SidebarHeader from '$lib/components/sidebar/SidebarHeader.svelte';
    import PackageListView from '$lib/components/package/PackageListView.svelte';
    import PackageFilter from '$lib/components/package/PackageFilter.svelte';

    const selectedProjectStore = getSelectedProjectStore();
    const toastStore = getToastStore();

    let selectedFilterType: number = 0;
    let searchQuery: string = "";
    let filterInstalled: boolean = false;
    const filterRadioOptions: RadioOption[] = [
        { value: 0, key: 'all', display: $t('home.sidebar.filter.option.all') },
        { value: 1, key: 'build', display: $t('home.sidebar.filter.option.build') },
        { value: 2, key: 'addon', display: $t('home.sidebar.filter.option.addon') },
    ];

    let fetchPackagesPromise: Promise<packageservice.ListPackagesResponse | undefined>;
    let dependencies: string[] = [];

    $: if ($selectedProjectStore) {
        loadDependencies();
    }

    async function loadDependencies() {
        var id = selectedProjectStore.latest();
        if (!id) {
            return;
        }

        var result = await GetProject(id);
        dependencies = result.project?.addons || [];
        dependencies.push(result.project?.build || '');
    }

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

        const addPackageToast: ToastSettings = {
            message: "Added Package!",
        };

        toastStore.trigger(addPackageToast);
    }

    function handleRefreshPackages(): void {
        const refreshPackageToast: ToastSettings = {
            message: "Refreshed Packages!",
        };

        toastStore.trigger(refreshPackageToast);
    }

    function handlePackageDownload(event: CustomEvent<{ packageId: string }>) {
        console.log('Downloaded package', event.detail.packageId);

        const downloadPackageToast: ToastSettings = {
            message: `Downloaded ${event.detail.packageId}!`,
        };

        toastStore.trigger(downloadPackageToast);
    }

    function handlePackageCancel(event: CustomEvent<{ packageId: string }>) {
        console.log('Cancel package', event.detail.packageId);
    }

    function handlePackageDelete(event: CustomEvent<{ packageId: string }>) {
        const deletePackageToast: ToastSettings = {
            message: `Deleted Package: ${event.detail.packageId}`,
        };

        toastStore.trigger(deletePackageToast);
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
    <PackageFilter
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
                <PackageListView
                    on:download={handlePackageDownload}
                    on:cancel={handlePackageCancel}
                    on:delete={handlePackageDelete}
                    packages={response.packages}
                    bind:dependencies={dependencies}    
                />
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