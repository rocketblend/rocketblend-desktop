<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    
    import { getToastStore, type ToastSettings } from '@skeletonlabs/skeleton';

    import { EventsOn } from '$lib/wailsjs/runtime';
    import { application, enums, types } from '$lib/wailsjs/go/models';
    import { ListPackages } from '$lib/wailsjs/go/application/Driver';

    import { t } from '$lib/translations/translations';
    import { createPackageStore } from '$lib/stores';
    import { EVENT_DEBOUNCE, SEARCH_STORE_INSERT_CHANNEL } from '$lib/events';
    import { debounce } from '$lib/utils';
    import type { RadioOption } from '$lib/types';

    import { PackageFilter, PackageList } from './package';
    import SidebarHeader from './sidebar-header.svelte';

    const packageStore = createPackageStore();
    const toastStore = getToastStore();

    const defaultFilterRadioOptions: RadioOption[] = [
        { value: undefined, display: $t('home.sidebar.filter.option.all') },
        { value: enums.PackageType.BUILD, display: $t('home.sidebar.filter.option.build') },
        { value: enums.PackageType.ADDON, display: $t('home.sidebar.filter.option.addon') },
    ];

    const fetchPackagesDebounced = debounce(fetchPackages, EVENT_DEBOUNCE);

    export let projectId: string | undefined;
    export let dependencies: string[];
    export let addonFeature: boolean;

    let selectedFilterType: enums.PackageType | undefined;
    let searchQuery: string = "";
    let filterInstalled: boolean;
    let filterRadioOptions: RadioOption[] = [];

    let initialLoad: boolean = true;
    let error: boolean = false;
    let cancelListener: () => void;

    function handleInputChange(): void {
        fetchPackages();
    }

    function handleAddPackage(): void {
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

    function fetchPackages() {
        error = false;

        let state = "";
        if (filterInstalled) {
            state = enums.PackageState.INSTALLED;
        }

        const opts = application.ListPackagesOpts.createFrom({
            query: searchQuery,
            type: selectedFilterType,
            state: state,
        });

        ListPackages(opts).then(result => {
            initialLoad = false;
            packageStore.set([...result.packages || []]);
        }).catch(error => {
            error = true;
            packageStore.set([]);
        });
    }
  
    onMount(() => {
        fetchPackages();

        cancelListener = EventsOn(SEARCH_STORE_INSERT_CHANNEL, (data: { id: string, indexType: string }) => {
            if (data.indexType === "package") {
                fetchPackagesDebounced();
            }
        });
    });

    onDestroy(() => {
        if (cancelListener) {
            cancelListener();
        }
    });

    $: {
        filterRadioOptions = addonFeature ? defaultFilterRadioOptions : [];
        selectedFilterType = addonFeature ? undefined : enums.PackageType.BUILD;
        fetchPackages();
    }
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
        {#if $packageStore && $packageStore.length > 0}
            <PackageList
                projectId={projectId}
                packages={$packageStore}
                dependencies={dependencies}    
            />
        {:else}
            {#if initialLoad}
                <div class="space-y-4 p-2">
                    {#each Array(10) as _}
                        <div class="placeholder animate-pulse p-5 h-10" />
                    {/each}
                </div>
            {:else if error}
                <p>{$t('home.sidebar.error')}</p>
            {:else}
                <div class="flex flex-col items-center justify-center p-2 gap-2">
                    <p class="font-bold text-sm text-surface-200 text-center">{$t('home.sidebar.noresults')}</p>
                </div>
            {/if}
        {/if}
    </div>
</div>