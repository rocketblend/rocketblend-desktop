<script lang="ts">
    import type { TableSource } from '@skeletonlabs/skeleton';
    import { Table } from '@skeletonlabs/skeleton';
  
    import type { project } from '$lib/wailsjs/go/models';
    import { tableMapperValues } from '$lib/components/core';
  
    export let sourceData: project.Project[];
  
    let tableSource: TableSource;
  
    // Whenever sourceData changes, this will run and update tableSource.
    $: {
      tableSource = {
        head: ['Project', 'file', 'Path', 'Tags'],
        body: tableMapperValues(sourceData, ['settings.name', 'blendFile.blendFileName', 'key','settings.tags']),
      };
    }
</script>
  
<Table source={tableSource} interactive={true} />