<script lang="ts">
    import type { types } from '$lib/wailsjs/go/models';

    import PackageListItem from './package-list-item.svelte';

    export let projectId: string | undefined;
    export let packages: types.Package[];
    export let dependencies: string[];
</script>

<div class="flex-auto space-y-1 rounded-token">
    {#each packages as pkg}
        <PackageListItem
            projectId={projectId}
            dependencies={dependencies}
            pack={pkg}
        />
    {/each}
</div>
