<script lang="ts">
    import "../app.postcss";
    
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';

    import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
    import { initializeStores, storePopup, getToastStore } from '@skeletonlabs/skeleton';
    import { Toast, AppBar, AppShell } from '@skeletonlabs/skeleton';

    import { Quit, WindowMinimise, WindowToggleMaximise } from '$lib/wailsjs/runtime';

    import { t } from '$lib/translations/translations';
    import { setupGlobalEventListeners, tearDownGlobalEventListeners } from '$lib/events';
    import { getLogStore } from "$lib/stores";

    import IconCloseFill from '~icons/ri/close-fill'
    import IconMoreFill from '~icons/ri/more-fill'
    import IconSubtractFill from '~icons/ri/subtract-fill'
    import IconCheckboxMultipleBlankLine from '~icons/ri/checkbox-multiple-blank-line'
    import IconHomeFill from '~icons/ri/home-fill'

    import Footer from "$lib/containers/Footer.svelte";
    import Sidebar from "$lib/containers/Sidebar.svelte";
	import UtilityDrawer from "$lib/containers/UtilityDrawer.svelte";

    initializeStores();

    const logStore = getLogStore();
    const toastStore = getToastStore();

    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

    function handleViewHome(): void {
        goto(`/`);
    }

    onMount(() => {
        setupGlobalEventListeners(logStore, toastStore);
    });

    onDestroy(() => {
        tearDownGlobalEventListeners();
    });

</script>

<UtilityDrawer/>

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
                <button type="button" class="btn btn-sm py-2 px-4 rounded-none text-2xl">
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