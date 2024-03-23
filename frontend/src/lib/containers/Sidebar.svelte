<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    
    import type {  ToastSettings } from '@skeletonlabs/skeleton';
    import { getToastStore } from '@skeletonlabs/skeleton';

    import { t } from '$lib/translations/translations';

    import { EventsOn } from '$lib/wailsjs/runtime';
    import { pack } from '$lib/wailsjs/go/models';
    import { GetProject, ListPackages, InstallPackageOperation } from '$lib/wailsjs/go/application/Driver';

    import type { RadioOption } from '$lib/types';
    import { EVENT_DEBOUNCE, SEARCH_STORE_INSERT_CHANNEL } from '$lib/events';
    import { getSelectedProjectStore, createPackageStore } from '$lib/stores';
    import { debounce } from '$lib/utils';

    import { SidebarHeader } from '$lib/components/ui/sidebar';
    import PackageListView from '$lib/components/ui/package/PackageListView.svelte';
    import PackageFilter from '$lib/components/ui/package/PackageFilter.svelte';

    const packageStore = createPackageStore();
    const selectedProjectStore = getSelectedProjectStore();
    const toastStore = getToastStore();

    const filterRadioOptions: RadioOption[] = [
        { value: pack.PackageType.UNKNOWN, display: $t('home.sidebar.filter.option.all') },
        { value: pack.PackageType.BUILD, display: $t('home.sidebar.filter.option.build') },
        { value: pack.PackageType.ADDON, display: $t('home.sidebar.filter.option.addon') },
    ];

    const fetchPackagesDebounced = debounce(fetchPackages, EVENT_DEBOUNCE);

    let selectedFilterType: number = 0;
    let searchQuery: string = "";
    let filterInstalled: boolean = false;
    let dependencies: string[] = [];

    let initialLoad: boolean = true;
    let error: boolean = false;
    let cancelListener: () => void;

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

    function handlePackageDownload(event: CustomEvent<{ packageId: string }>) {
        const packageId = event.detail.packageId;

        InstallPackageOperation(packageId).then((result) => {
            const downloadPackageToast: ToastSettings = {
                message: `Downloading Package: ${packageId}`,
            };

            toastStore.trigger(downloadPackageToast);
        }).catch(error => {
            const downloadPackageToast: ToastSettings = {
                message: `Error starting package download: ${error}`,
                background: "variant-filled-error"
            };

            toastStore.trigger(downloadPackageToast);
        });
    }

    function handlePackageCancel(event: CustomEvent<{ packageId: string }>) {
        const packageId = event.detail.packageId;
        //var item = $packageStore.find((pack) => pack.id?.toString() === packageId);

        const cancelledPackageToast = {
            message: `Cancelled ${packageId}!`,
        };
        toastStore.trigger(cancelledPackageToast);
    }

    function handlePackageDelete(event: CustomEvent<{ packageId: string }>) {
        const deletePackageToast: ToastSettings = {
            message: `Deleted Package: ${event.detail.packageId}`,
        };

        toastStore.trigger(deletePackageToast);
    }

    function fetchPackages() {
        error = false;
        ListPackages(searchQuery, selectedFilterType, filterInstalled).then(result => {
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
            <PackageListView
                on:download={handlePackageDownload}
                on:cancel={handlePackageCancel}
                on:delete={handlePackageDelete}
                packages={$packageStore}
                bind:dependencies={dependencies}    
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
                <div class="p-2">
                    <p class="font-bold text-sm text-surface-200 text-center">{$t('home.sidebar.noresults')}</p>
                </div>
            {/if}
        {/if}
    </div>
</div>