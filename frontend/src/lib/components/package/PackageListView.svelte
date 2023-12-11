<script lang="ts">
    import { ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
    import PackageListItem from './PackageListItem.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    export let packages: pack.Package[] = [];
    export let selectedPackageIds: string[] = [];
</script>

<ListBox class="flex-auto" multiple>
    {#each packages as pack}
        <ListBoxItem bind:group={selectedPackageIds} name="packages" value={pack.id} active="variant-glass-primary" hover="hover:variant-filled-surface" rounded="rounded" class="truncate overflow-hidden">
            <PackageListItem name={pack.name || ""} reference={pack.reference || ""} />
        </ListBoxItem>
    {/each}
</ListBox>
