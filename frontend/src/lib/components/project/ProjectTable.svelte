<script lang="ts">
    import type { project } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core'

    import Table from '$lib/components/core/table/Table.svelte';
    import type { TableSource } from '$lib/components/core/table/types.js';

    export let sourceData: project.Project[];
    export let selectedProjectIds: string[] = [];

    let tableSource: TableSource;

    $: {
        tableSource = {
            head: ['Project', 'File', 'Build', 'Tags'],
            body: tableMapperValues(sourceData, 'id', ['name', 'fileName', 'build', 'tags']),
        };
    }
</script>

<Table bind:source={tableSource} interactive={true} bind:selected={selectedProjectIds} />