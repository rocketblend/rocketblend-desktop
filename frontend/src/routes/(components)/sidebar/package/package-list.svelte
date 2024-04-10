<script lang="ts">
    import type { pack } from '$lib/wailsjs/go/models';

    import PackageListItem from './package-list-item.svelte';

    export let packages: pack.Package[] = [];
    export let dependencies: string[] = [];
</script>

<div class="flex-auto space-y-1 rounded-token">
    {#each packages as pkg}
        <PackageListItem
            dependencies={dependencies}
            pack={pkg}
        />
    {/each}
</div>
