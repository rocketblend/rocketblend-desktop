<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    import type { project } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core'

    import Table from '$lib/components/core/table/Table.svelte';
    import type { TableSource } from '$lib/components/core/table/types.js';
  
    export let sourceData: project.Project[];
  
    let tableSource: TableSource;

    const dispatch = createEventDispatcher<{ selected: project.Project | null }>();

    function handleSelected(event: CustomEvent<string[]>) {
      var project = sourceData.find((p) => p.id?.toString() === event.detail[0]);
      dispatch('selected', project);
    }
  
    $: {
      tableSource = {
        head: ['Project', 'File', 'Build', 'Tags'],
        body: tableMapperValues(sourceData, ['name', 'fileName', 'build', 'tags']),
        meta: tableMapperValues(sourceData, ['id']),
      };
    }
</script>
  
<Table source={tableSource} interactive={true} on:selected={handleSelected} />