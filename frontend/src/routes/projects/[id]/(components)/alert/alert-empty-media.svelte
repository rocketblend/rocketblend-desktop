<script lang="ts">
    import {
        Alert,
        AlertTitle,
        AlertDescription,
    } from '$lib/components/ui/alert';

    export let folder = "media";
</script>
  
<Alert>
    <svelte:fragment slot="title">
        <AlertTitle >
            No media files found
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        You can images, videos, and other media files to this project. Just create a <code class="code">{folder}</code> folder in the root of your project and add your files there. 
    </AlertDescription>
</Alert>