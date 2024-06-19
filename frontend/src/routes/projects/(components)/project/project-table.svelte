<script lang="ts">
    import { type types, enums } from '$lib/wailsjs/go/models';

    import { Table, tableMapperValues } from '$lib/components/ui/table'
    import type { TableSource, TableColumn } from '$lib/components/ui/table'

    export let sourceData: types.Project[];
    export let selectedProjectIds: string[] = [];

    let tableSource: TableSource;

    $: data = sourceData.map((project) => ({
        id: project.id.toString(),
        name: project.name || "",
        fileName: project.fileName || "",
        build: project.dependencies.find((d) => d.type === enums.PackageType.BUILD)?.reference || "",
    }));

    $: {
        tableSource = {
            head: [
                { label: 'name', display: 'Name', sortable: true },
                { label: 'filename', display: 'File Name', sortable: true },
                { label: 'build', display: 'Build', sortable: true },
            ] as TableColumn[],
            body: tableMapperValues(data, 'id', ['name', 'fileName', 'build']),
            foot: ['']
        };
    }
</script>

<Table source={tableSource} interactive={true} bind:selected={selectedProjectIds} on:dblclick/>