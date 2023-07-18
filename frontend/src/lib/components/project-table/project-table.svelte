<script lang="ts">
    import { Table, tableMapperValues } from '@skeletonlabs/skeleton';
    import type { TableSource } from '@skeletonlabs/skeleton';

    import type { project } from '$lib/wailsjs/go/models';

    export let sourceData: project.Project[];

    const tableSource: TableSource = {
        head: ['Project', 'Key', 'Last Opened'],
        body: tableMapperValues(sourceData, ['name', 'key', 'lastUpdated']),
        meta: tableMapperValues(sourceData, ['name', 'key', 'lastUpdated']),
    };
</script>

<Table source={tableSource} interactive={true} />
