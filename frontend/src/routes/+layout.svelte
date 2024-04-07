<script lang="ts">
    import "../app.postcss";
    
    import type { LayoutData } from "./$types";
    import { onMount, onDestroy } from 'svelte';

    import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
    import { initializeStores, storePopup, getToastStore, getDrawerStore } from '@skeletonlabs/skeleton';
    import { Toast, AppBar, AppShell, Modal } from '@skeletonlabs/skeleton';

    import { Quit, WindowMinimise, WindowToggleMaximise } from '$lib/wailsjs/runtime';

    import { t } from '$lib/translations/translations';
    import { setupGlobalEventListeners, tearDownGlobalEventListeners } from '$lib/events';
    import { getLogStore } from "$lib/stores";

    import IconCloseFill from '~icons/ri/close-fill'
    import IconMoreFill from '~icons/ri/more-fill'
    import IconSubtractFill from '~icons/ri/subtract-fill'
    import IconCheckboxMultipleBlankLine from '~icons/ri/checkbox-multiple-blank-line'
    import IconHomeFill from '~icons/ri/home-fill'
    import IconTerminalBoxFill from '~icons/ri/terminal-box-fill';
    import IconBrainFill from '~icons/ri/brain-fill'
    import IconSettingsFill from '~icons/ri/settings-4-fill'

    import IconArrowLeftFile from '~icons/ri/arrow-left-s-line'
    import IconArrowRightFile from '~icons/ri/arrow-right-s-line'

    import { Footer, Sidebar, UtilityDrawer, Breadcrumb } from "./(components)"

    import Logo from "$lib/assets/images/logo-slim.png?enhanced"

    initializeStores();

    const logStore = getLogStore();
    const toastStore = getToastStore();
    const drawerStore = getDrawerStore();

    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

    export let data: LayoutData;

    // Function to navigate back
    function goBack() {
        window.history.back();
    }

    // Function to navigate forward
    function goForward() {
        window.history.forward();
    }

    function openTerminal() {
        drawerStore.open();
    }

    onMount(() => {
        setupGlobalEventListeners(logStore, toastStore);
    });

    onDestroy(() => {
        tearDownGlobalEventListeners();
    });
</script>

<Modal
    rounded="rounded"
/>

<UtilityDrawer developer={data.preferences.feature.developer}/>

<Toast
    zIndex="z-40"
    background="variant-filled-surface"
    padding="p-4"
    position="br"
    rounded="rounded"
    class="mx-4 mt-10 mb-24"
/>

<AppShell slotSidebarLeft="flex flex-col overflow-y-hidden space-y-2 pl-2 w-96 h-full" slotPageContent="overflow-hidden h-full">
    <svelte:fragment slot="header">
        <div style="--wails-draggable:drag">
            {#if data.details.platform !== "macos/intel" && data.details.platform !== "macos/apple" }
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
            {:else}
                <div class="h-9 p-0"></div>
            {/if}
        </div>
    </svelte:fragment>
    <svelte:fragment slot="sidebarLeft" >
        <div class="h-32">
            <div class="relative rounded-container-token overflow-hidden h-full">
                <div class="absolute w-full h-full">
                    <a class="flex items-center h-full px-4 gap-2" href="/">
                        <div>
                            <span class="h4 font-bold">RocketBlend</span><br>
                            <span class="h5 text-surface-800-100-token">Desktop</span><br>
                            <span class="text-sm text-surface-500-400-token">{data.details.version}</span>
                        </div>
                    </a>
                </div>
                <enhanced:img src={Logo} alt=""/>
            </div>
        </div>
        <div class="card flex-grow shadow-none p-4 overflow-hidden">
            <Sidebar addon={data.preferences.feature.addon}/>
        </div>
    </svelte:fragment>
    <svelte:fragment slot="pageHeader">
        <div class="h-full p-2 pt-0">
            <div class="shadow-none card px-6 py-4 h-full">
                <div class="flex flex-wrap justify-between items-center gap-6">
                    <div class="flex justify-between items-center gap-6">
                        <div>
                            <button type="button" class="btn btn-sm variant-filled-surface" on:click={goBack}><IconArrowLeftFile/></button>
                            <button type="button" class="btn btn-sm variant-filled-surface" on:click={goForward}><IconArrowRightFile/></button>
                            <a class="btn btn-sm variant-filled-surface" href="/"><IconHomeFill/></a>
                        </div>
                        <div class="flex">
                            <Breadcrumb />
                        </div>
                    </div>

                    <div class="flex items-center gap-4">
                        <button type="button" class="btn text-lg text-surface-700-200-token p-1" on:click={openTerminal}>
                            <IconTerminalBoxFill/>
                        </button>
                        <!-- <a class="btn text-lg text-surface-700-200-token px-2" href="/metrics">
                            <IconBrainFill/>
                        </a> -->
                        <a class="btn text-lg text-surface-700-200-token p-1" href="/preferences">
                            <IconSettingsFill/>
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </svelte:fragment>
    <div class="h-full p-2 py-0">
        <div class="shadow-none card p-6 h-full overflow-hidden">
            <slot />
        </div>
    </div>
    <svelte:fragment slot="footer">
        <Footer/>
    </svelte:fragment>
</AppShell>