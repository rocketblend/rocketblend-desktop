<script lang="ts">
    import { OpenExplorer } from '$lib/wailsjs/go/application/Driver';
    import { application } from '$lib/wailsjs/go/models';

    import { AlertAction } from '$lib/components/ui/alert';

    export let path: string;
    
    let text: string = "View Location";

    let disabled = false && !path;

    async function explore() {
        if (disabled || !path) {
            return;
        }

        disabled = true;
        const opts = application.OpenExplorerOptions.createFrom({
            path: path,
        });

        await OpenExplorer(opts);

        disabled = false;
    }
</script>

<AlertAction on:click={explore} disabled={disabled}>
    {text}
</AlertAction>
