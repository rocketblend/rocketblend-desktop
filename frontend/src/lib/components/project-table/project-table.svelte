<script lang="ts">
    import type { TableSource } from '@skeletonlabs/skeleton';
    import { Table } from '@skeletonlabs/skeleton';

    import type { project } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core';

    export let sourceData: project.Project[];

    if (!sourceData) {
        sourceData = [];
    }

    const tableSource: TableSource = {
        head: ['Project', 'file', 'Path', 'Tags'],
        body: tableMapperValues(sourceData, ['settings.name', 'blendFile.blendFileName', 'key','settings.tags']),
    };
</script>

{#if sourceData.length > 0}
    <Table source={tableSource} interactive={true} />
{:else}
    <p>No projects found!</p>
{/if}
