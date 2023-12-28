<script lang="ts">
    import "../app.postcss";
    
    import { onMount, onDestroy } from 'svelte';
    import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
    import { initializeStores, getToastStore, getDrawerStore, storePopup } from '@skeletonlabs/skeleton';
    import { Toast, AppBar, AppShell, Drawer } from '@skeletonlabs/skeleton';
    import { TabGroup, Tab, TabAnchor } from '@skeletonlabs/skeleton';
    import type { ToastSettings } from '@skeletonlabs/skeleton';

    import { goto } from '$app/navigation';

    import IconCloseFill from '~icons/ri/close-fill'
    import IconMoreFill from '~icons/ri/more-fill'
    import IconSubtractFill from '~icons/ri/subtract-fill'
    import IconCheckboxMultipleBlankLine from '~icons/ri/checkbox-multiple-blank-line'
    import IconHomeFill from '~icons/ri/home-fill'

    import { t } from '$lib/translations/translations';
    import { EventsEmit, EventsOff, EventsOn, Quit, WindowMinimise, WindowToggleMaximise } from '$lib/wailsjs/runtime';

    import Footer from "$lib/containers/Footer.svelte";
    import Sidebar from "$lib/containers/Sidebar.svelte";
	import type { LogEvent } from "$lib/types";
	import LogFeed from "$lib/components/feed/LogFeed.svelte";

    initializeStores();
    const toastStore = getToastStore();
    const drawerStore = getDrawerStore();

    const MAX_LOGS= 250;

    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

    let drawTabSet: number = 0;
    let logs: LogEvent[] = [];

    function handleViewHome(): void {
        goto(`/`);
    }

    function handleToggleMoreClick() {
        drawerStore.open();
    }

    function addLog(newLog: LogEvent) {
        if (logs.length >= MAX_LOGS) {
            logs.shift();
        }

        logs = [...logs, newLog];
    }

    onMount(() => {     
        EventsOn('launchArgs', (data: { args: string[]}) => {
            console.log('launchArgs', data)
            if (data.args && data.args.length !== 0) {
                    var launchToast: ToastSettings = {
                    message: `Args: ${data.args.join(', ')}`,
                    timeout: 5000,
                };

                toastStore.trigger(launchToast);
            }
        });

        EventsOn('logStream', (data: LogEvent) => {
            addLog(data);
        });

        const launchToast: ToastSettings = {
            message: t.get('home.greeting'),
            background: 'bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500 text-white',
            timeout: 5000,
        };

        toastStore.trigger(launchToast);

        EventsEmit('ready'); // Notify Wails that the frontend is ready for events
    });

    onDestroy(() => {
        EventsOff('launchArgs');
        EventsOff('logStream');
    });
</script>

<Drawer
    class="h-full overflow-hidden"
    position="bottom"
    rounded="none"
    zIndex="z-50">
        <TabGroup class="flex flex-col h-full overflow-hidden" active="border-b-2 border-primary-400-500-token" rounded="" regionPanel="px-4 pb-4 flex-grow overflow-hidden h-full" regionList="flex flex-none">
            <Tab bind:group={drawTabSet} name="tab1" value={0}>{$t('home.drawer.tab.output.header')}</Tab>
            <Tab bind:group={drawTabSet} name="tab2" value={1}>{$t('home.drawer.tab.console.header')}</Tab>
            <svelte:fragment slot="panel">
                {#if drawTabSet === 0}
                    <LogFeed bind:feed={logs}/>
                {:else if drawTabSet === 1}
                    <div class="flex justify-center items-center h-full">
                        <p>{$t('home.soon')}</p>
                    </div>
                {/if}
            </svelte:fragment>
        </TabGroup>
</Drawer>

<Toast
    zIndex="z-40"
    background="variant-filled-surface"
    padding="p-4"
    position="br"
    rounded="rounded"
    class="mx-4 mt-10 mb-24"
/>

<AppShell slotSidebarLeft="flex flex-col overflow-y-hidden space-y-2 pl-2 w-96 h-full">
    <svelte:fragment slot="header">
        <div style="--wails-draggable:drag">
            <AppBar background="bg-surface-50-900-token" padding="p0" slotTrail="space-x-0 -mt-3">
                <svelte:fragment slot="lead">
                <button type="button" class="btn btn-sm py-2 px-4 rounded-none text-2xl" on:click={handleToggleMoreClick} >
                    <IconMoreFill/>
                </button>
                </svelte:fragment>
                <svelte:fragment slot="trail">
                <button type="button" class="btn btn-sm py-2 px-4 hover:bg-stone-700 rounded-none" on:click={WindowMinimise}>
                    <IconSubtractFill/>
                </button>
                <button type="button" class="btn btn-sm py-2 px-4 hover:bg-stone-700 rounded-none" on:click={WindowToggleMaximise}>
                    <IconCheckboxMultipleBlankLine/>
                </button>
                <button type="button" class="btn btn-sm py-2 px-4 hover:bg-red-700 rounded-none" on:click={Quit}>
                    <IconCloseFill/>
                </button>
                </svelte:fragment>
            </AppBar>
        </div>
    </svelte:fragment>
    <svelte:fragment slot="sidebarLeft" >
        <div class="card flex-shrink-0 p-4 shadow-none">
            <button type="button" class="btn btn-sm py-2 px-4 pl-0 text-lg text-surface-200" on:click={handleViewHome}>
                <IconHomeFill/>
                <span class="font-bold">{$t('home.navigation.root')}</span>
            </button>
        </div>
        <div class="card flex-grow shadow-none p-4 overflow-hidden">
            <Sidebar/>
        </div>
    </svelte:fragment>
    <div class="h-full p-2 py-0">
        <div class="shadow-none card p-6 h-full">
            <slot />
        </div>
    </div>
    <svelte:fragment slot="footer">
        <Footer/>
    </svelte:fragment>
</AppShell>