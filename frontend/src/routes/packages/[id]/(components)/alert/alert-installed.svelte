<script lang="ts">
    import IconCheckFill from '~icons/ri/check-fill';
    import { Alert, AlertTitle, AlertDescription, AlertAction } from '$lib/components/ui/alert';
	import { ActionExplore } from '../action';

    export let path: string;
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
        <!-- <AlertAction text="Delete" disabled/> -->
        <ActionExplore path={path} />
    </svelte:fragment>
</Alert>