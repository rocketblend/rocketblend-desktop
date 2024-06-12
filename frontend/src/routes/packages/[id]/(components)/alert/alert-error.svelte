<script lang="ts">
    import IconWarningFill from '~icons/ri/file-damage-fill'
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconWarningFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Error
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription message="An error occurred while downloading the package."/>
    <svelte:fragment slot="actions">
        <AlertAction disabled>
            Retry
        </AlertAction>
    </svelte:fragment>
</Alert>