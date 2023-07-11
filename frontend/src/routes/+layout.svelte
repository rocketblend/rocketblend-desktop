<script>
  import "../theme.postcss";
  import "@skeletonlabs/skeleton/styles/all.css";
  import "../app.postcss";

  import { AppBar, AppShell } from '@skeletonlabs/skeleton';

  import IconCloseFill from '~icons/ri/close-fill'
  import IconMoreFill from '~icons/ri/more-fill'
  import IconSubtractFill from '~icons/ri/subtract-fill'
  import IconCheckboxMultipleBlankLine from '~icons/ri/checkbox-multiple-blank-line'

  import { Quit, WindowMinimise, WindowToggleMaximise } from '$lib/wailsjs/runtime';
</script>

<AppShell>
  <svelte:fragment slot="header">
    <div style="--wails-draggable:drag">
      <AppBar background="bg-surface-100-900-token" padding="p0" slotTrail="space-x-0 -mt-3">
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
  <div class="h-full p-2 pt-0">
    <div class="shadow-none card p-4 h-full">
      <slot />
    </div>
  </div>
</AppShell>
