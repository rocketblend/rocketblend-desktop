<script lang="ts">
    import type { TableSource } from '@skeletonlabs/skeleton';
    import { Table } from '@skeletonlabs/skeleton';
  
    import type { projectservice } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core';
  
    export let sourceData: projectservice.Project[];
  
    let tableSource: TableSource;
  
    // Whenever sourceData changes, this will run and update tableSource.
    $: {
      tableSource = {
        head: ['Project', 'File', 'Build', 'Tags'],
        body: tableMapperValues(sourceData, ['name', 'fileName', 'build', 'tags']),
      };
    }
</script>
  
<Table source={tableSource} interactive={true} />