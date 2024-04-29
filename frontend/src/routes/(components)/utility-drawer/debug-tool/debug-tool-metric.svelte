<script lang="ts">
    import { AggregateMetrics } from '$lib/wailsjs/go/application/Driver';
    import { application, types } from '$lib/wailsjs/go/models';

    let domain: string = "";
    let name: string= "";;
    let aggregate: types.Aggregate;

    const fetchMetrics = async () => {
        const request = application.AggregateMetricsOpts.createFrom({
            domain: domain,
            name: name,
        });

        AggregateMetrics(request).then((response) => {
            aggregate = response;
        }).catch((error) => {
            console.error(error);
        });
    };

</script>

<div class="flex flex-col card p-2 space-y-2">
    <input type="text" class="input" bind:value={domain} placeholder="Domain" />
    <input type="text" class="input" bind:value={name} placeholder="Name" />
    <button class="btn variant-filled" on:click={fetchMetrics}>Aggregate</button>
    <hr>
    <div>
        <pre>{JSON.stringify(aggregate, null, 2)}</pre>
    </div>
</div>