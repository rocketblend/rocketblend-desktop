<script>
    import "../theme.postcss";
    import "@skeletonlabs/skeleton/styles/all.css";
    import "../app.postcss";

    import { AppBar, AppShell } from '@skeletonlabs/skeleton';
    import { Quit, WindowMinimise, WindowToggleMaximise } from '$lib/wailsjs/runtime';

    import IconCloseFill from '~icons/ri/close-fill'
    import IconMoreFill from '~icons/ri/more-fill'
    import IconSubtractFill from '~icons/ri/subtract-fill'
    import IconCheckboxMultipleBlankLine from '~icons/ri/checkbox-multiple-blank-line'

    import Footer from '$lib/components/footer/footer.svelte';
	  import Sidebar from "$lib/components/sidebar/sidebar.svelte";
</script>

<AppShell>
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
  <svelte:fragment slot="sidebarLeft">
    <div class="h-full p-2 pr-0 py-0 w-80">
      <div class="shadow-none card p-3 pt-6 h-full overflow-y-auto overflow-x-hidden">
        <Sidebar/>
      </div>
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