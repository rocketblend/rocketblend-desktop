<script lang="ts">
    import IconCheckFill from '~icons/ri/check-fill';
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
	import { ActionExplore, ActionRemove } from '../action';

    export let path: string;
    export let packageId: string;
</script>

<Alert>
    <svelte:fragment slot="icon">
        <IconCheckFill class="text-2xl"/>
    </svelte:fragment>
    <svelte:fragment slot="title">
        <AlertTitle >
            Installed
        </AlertTitle>
    </svelte:fragment>
    <AlertDescription message="Package is ready to be used."/>
    <svelte:fragment slot="actions">
        <ActionRemove packageId={packageId}/>
        <ActionExplore path={path} />
    </svelte:fragment>
</Alert>