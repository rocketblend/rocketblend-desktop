<script lang="ts">
    import {
        Alert,
        AlertTitle,
        AlertDescription,
    } from '$lib/components/ui/alert';

    import IconImage from '~icons/ri/image-add-fill';

    export let folder = "media";
</script>
  
<Alert>
    <svelte:fragment slot="icon">
        <IconImage class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Personalise your project
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        You can display images, videos, and other media files on a project. Just create a <code class="code">{folder}</code> folder in the root of this project and add your files there. 
    </AlertDescription>
</Alert>