<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import PackageListItem from './package-list-item.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    const dispatch = createEventDispatcher();

    export let packages: pack.Package[] = [];
    export let dependencies: string[] = [];

    function handleItemClick(reference: string = "") {
        console.log("click", reference);

        if (!dependencies.includes(reference)) {
            dependencies = [...dependencies, reference];
            return
        }
        
        dependencies = dependencies.filter(dep => dep !== reference);
    }

    function handleItemDownload(packageId: string = "") {
        dispatch('download', { packageId });
    }

    function handleItemCancel(packageId: string = "") {
        dispatch('cancel', { packageId });
    }

    function handleItemDelete(packageId: string = "") {
        dispatch('delete', { packageId });
    }
</script>

<div class="flex-auto space-y-1 rounded-token">
    {#each packages as pkg}
        <PackageListItem
            on:click={() => handleItemClick(pkg.reference?.toString())}
            on:download={() => handleItemDownload(pkg.id?.toString())}
            on:cancel={() => handleItemCancel(pkg.id?.toString())}
            on:delete={() => handleItemDelete(pkg.id?.toString())}
            selected={dependencies.includes(pkg.reference?.toString() || "")}
            name={pkg.name?.toString() || ""}
            tag={pkg.tag?.toString() || ""}
            verified={pkg.verified? true : false}
            reference={pkg.reference?.toString() || ""}
            platform={pkg.platform?.toString() || ""}
            version={pkg.version?.toString() || ""}
            author={pkg.author?.toString() || ""}
            type={pkg.type}
            state={pkg.state}
        />
    {/each}
</div>
