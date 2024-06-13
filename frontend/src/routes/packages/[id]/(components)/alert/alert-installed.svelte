<script lang="ts">
    import IconCheckFill from '~icons/ri/check-fill';
    import { Alert, AlertTitle, AlertDescription } from '$lib/components/ui/alert';
	import { ActionExplore, ActionRemove } from '../action';

    export let path: string;
    export let packageId: string;

    $: title = path ? "Installed" : "Pre-Installed";
    $: description = path ? "Package is ready to be used." : "Package is ready to be used, as it comes pre-installed with blender.";
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconCheckFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            {title}
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription>
       {description}
    </AlertDescription>
    <svelte:fragment slot="actions">
        {#if path}
            <ActionRemove packageId={packageId}/>
            <ActionExplore path={path} />
        {/if}
    </svelte:fragment>
</Alert>