<script lang="ts">
    import IconErrorWarningFill from '~icons/ri/error-warning-fill'

    import {
        Alert,
        AlertTitle,
        AlertDescription,
    } from '$lib/components/ui/alert';
</script>
  
<Alert variant="ghost-error">
    <svelte:fragment slot="icon">
        <IconErrorWarningFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Build Required
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
        Please add a build dependency to use the project.
    </AlertDescription>
</Alert>