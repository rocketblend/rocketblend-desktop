<script lang="ts">
    import type { project } from '$lib/wailsjs/go/models';

    import { Table, tableMapperValues } from '$lib/components/ui/table'
    import type { TableSource, TableColumn } from '$lib/components/ui/table'

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
            foot: ['']
        };
    }
</script>

<Table source={tableSource} interactive={true} bind:selected={selectedProjectIds} on:sortChanged on:itemDoubleClick/>