<script lang="ts">
    import { t } from '$lib/translations/translations';
    import { Drawer, TabGroup, Tab } from '@skeletonlabs/skeleton';

    import { LogFeed } from "./log-feed"
    import {
        DebugToolMetrics,
        DebugToolOperations,
        DebugToolRequests
    } from "./debug-tool";

    export let developer: boolean = false;

    let drawTabSet: number = 0;

    $: {
        if (!developer) {
            drawTabSet = 0;
        }
    }
</script>

<Drawer
    class="h-full overflow-hidden"
    position="bottom"
    rounded="none"
    zIndex="z-50">
    <TabGroup class="flex flex-col h-full overflow-hidden" active="border-b-2 border-primary-400-500-token" rounded="" regionPanel="px-4 pb-4 flex-grow overflow-hidden h-full" regionList="flex flex-none">
        <Tab bind:group={drawTabSet} name="tab1" value={0}>{$t('home.drawer.tab.output')}</Tab>
        {#if developer}
            <Tab bind:group={drawTabSet} name="tab2" value={1}>{$t('home.drawer.tab.terminal')}</Tab>
            <Tab bind:group={drawTabSet} name="tab3" value={2}>{$t('home.drawer.tab.debug')}</Tab>
        {/if}
        <svelte:fragment slot="panel">
            <div class="h-full" hidden={!(drawTabSet == 0)}>
                <LogFeed />
            </div>
            <div class="h-full" hidden={!(drawTabSet == 1)}>
                <div class="flex justify-center items-center h-full">
                    <p>{$t('home.soon')}</p>
                </div>
            </div>
            <div class="h-full" hidden={!(drawTabSet == 2)}>
                <div class="overflow-auto grid grid-cols-2 md:grid-cols-3 gap-4 h-full">
                    <div>
                        <DebugToolRequests/>
                    </div>
                    <div>
                        <DebugToolOperations/>
                    </div>
                    <div>
                        <DebugToolMetrics/>
                    </div>
                </div>
            </div>
        </svelte:fragment>
    </TabGroup>
</Drawer>