<script lang="ts">
    import IconCloseFill from '~icons/ri/close-fill';
    import {
        Alert,
        AlertTitle,
        AlertDescription,
        AlertAction
    } from '$lib/components/ui/alert';
</script>
  
<Alert>
    <svelte:fragment slot="icon">
        <IconCloseFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle title="Disabled"/>
    </svelte:fragment>
    <AlertDescription message="Package is currently disabled on the selected project"/>
    <svelte:fragment slot="actions">
        <AlertAction text="Enable" on:click/>
    </svelte:fragment>
</Alert>