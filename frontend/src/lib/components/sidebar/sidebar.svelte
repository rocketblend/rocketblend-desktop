<script lang="ts">
    import { onMount } from 'svelte';
    import { t } from '$lib/translations/translations';

    import SearchInput from '$lib/components/core/search-input/search-input.svelte';

    import type { packageservice } from '$lib/wailsjs/go/models';
    import { ListPackages } from '$lib/wailsjs/go/application/Driver'

    let query: string = "";
    let fetchPackagesPromise : Promise<packageservice.ListPackagesResponse| undefined> ;

    async function fetchPackages(query:string): Promise<packageservice.ListPackagesResponse | undefined> {
        try {
            return await ListPackages(query);
        } catch (error) {
            console.error('Error fetching packages:', error);
            return undefined;
        }
    }

    onMount(() => {
        fetchPackagesPromise = fetchPackages(query);
    });

    function handleInputChange(event: Event): void {
        fetchPackagesPromise = fetchPackages(query);
    }

</script>

<div class="space-y-4">
    <h5 class="font-bold text-surface-200">{$t('home.sidebar.title')}</h5>
    <SearchInput bind:value={query} placeholder={$t('home.sidebar.search')} debounceDelay={500} on:input={handleInputChange}/>
    <div>
        {#await fetchPackagesPromise}
            <div class="flex-auto space-y-4 p-2">
                {#each Array(10) as _}
                    <div class="placeholder animate-pulse p-5" />
                {/each}
            </div>
        {:then response}
            {#if response && response.packages}
            <dl class="list-dl">
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
