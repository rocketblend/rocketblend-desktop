<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    import type { TableSource } from '@skeletonlabs/skeleton';
    import { Table } from '@skeletonlabs/skeleton';
  
    import type { projectservice } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core'
  
    export let sourceData: projectservice.Project[];
  
    let tableSource: TableSource;

    const dispatch = createEventDispatcher<{ selected: projectservice.Project | null }>();

    function handleSelected(event: CustomEvent<string[]>) {
      var project = sourceData.find((p) => p.id === event.detail[0]);
      dispatch('selected', project);
    }
  
    $: {
      tableSource = {
        head: ['Project', 'id', 'File', 'Build', 'Tags'],
        body: tableMapperValues(sourceData, ['name', 'id', 'fileName', 'build', 'tags']),
        meta: tableMapperValues(sourceData, ['id']),
      };
    }
</script>
  
<Table source={tableSource} interactive={true} on:selected={handleSelected} />