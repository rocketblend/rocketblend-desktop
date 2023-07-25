<script lang="ts">
    import { Table, tableMapperValues } from '@skeletonlabs/skeleton';
    import type { TableSource } from '@skeletonlabs/skeleton';

    import type { project } from '$lib/wailsjs/go/models';

    export let sourceData: project.Project[];

    if (!sourceData) {
        sourceData = [];
    }

    const tableSource: TableSource = {
        head: ['Project', 'Key'],
        body: tableMapperValues(sourceData, ['settings', 'key']),
    };
</script>

{#if sourceData.length > 0}
    <Table source={tableSource} interactive={true} />
{:else}
    <p>No projects found!</p>
{/if}
