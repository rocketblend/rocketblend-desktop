<script lang="ts">
    import type { project } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core'

    import Table from '$lib/components/core/table/Table.svelte';
    import type { TableSource, TableColumn } from '$lib/components/core/table/types.js';

    export let sourceData: project.Project[];
    export let selectedProjectIds: string[] = [];

    let tableSource: TableSource;

    $: {
        tableSource = {
            head: [
                { label: 'name', display: 'Name', sortable: true },
                { label: 'filename', display: 'File Name', sortable: true },
                { label: 'build', display: 'Build', sortable: true },
            ] as TableColumn[],
            body: tableMapperValues(sourceData, 'id', ['name', 'fileName', 'build']),
        };
    }
</script>

<Table source={tableSource} interactive={true} bind:selected={selectedProjectIds} on:sortChanged/>