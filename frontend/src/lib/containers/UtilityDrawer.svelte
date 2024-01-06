<script lang="ts">
    import { t } from '$lib/translations/translations';
    import { Drawer, TabGroup, Tab } from '@skeletonlabs/skeleton';

    import { getLogStore } from '$lib/stores';

    import RequestCancellationDebugTool from "$lib/containers/debug/RequestCancellationDebugTool.svelte";
    import LogFeed from "$lib/components/feed/LogFeed.svelte";
	import OperationDebugTool from './debug/OperationDebugTool.svelte';
	import MetricAggregateDebugTool from './debug/MetricAggregateDebugTool.svelte';

    const logStore = getLogStore();

    let drawTabSet: number = 0;
</script>

<Drawer
    class="h-full overflow-hidden"
    position="bottom"
    rounded="none"
    zIndex="z-50">
    <TabGroup class="flex flex-col h-full overflow-hidden" active="border-b-2 border-primary-400-500-token" rounded="" regionPanel="px-4 pb-4 flex-grow overflow-hidden h-full" regionList="flex flex-none">
        <Tab bind:group={drawTabSet} name="tab1" value={0}>{$t('home.drawer.tab.output')}</Tab>
        <Tab bind:group={drawTabSet} name="tab2" value={1}>{$t('home.drawer.tab.terminal')}</Tab>
        <Tab bind:group={drawTabSet} name="tab3" value={2}>{$t('home.drawer.tab.debug')}</Tab>
        <svelte:fragment slot="panel">
            <div class="h-full" hidden={!(drawTabSet == 0)}>
                <LogFeed feed={$logStore} />
            </div>
            <div class="h-full" hidden={!(drawTabSet == 1)}>
                <div class="flex justify-center items-center h-full">
                    <p>{$t('home.soon')}</p>
                </div>
            </div>
            <div class="h-full" hidden={!(drawTabSet == 2)}>
                <div class="overflow-auto grid grid-cols-2 md:grid-cols-3 gap-4 h-full">
                    <div>
                        <RequestCancellationDebugTool/>
                    </div>
                    <div>
                        <OperationDebugTool/>
                    </div>
                    <div>
                        <MetricAggregateDebugTool/>
                    </div>
                </div>
            </div>
        </svelte:fragment>
    </TabGroup>
</Drawer>