<script>
  import "../theme.postcss";
  import "@skeletonlabs/skeleton/styles/all.css";
  import "../app.postcss";

  import { AppBar, AppShell } from '@skeletonlabs/skeleton';

  import IconCloseRounded from '~icons/material-symbols/close-rounded'
  import IconSelectWindowSharp from '~icons/material-symbols/select-window-sharp'
  import IconRemove from '~icons/material-symbols/remove'
  import IconMoreHoriz from '~icons/material-symbols/more-horiz'
</script>

<AppShell>
  <svelte:fragment slot="header">
    <div style="--wails-draggable:drag">
      <AppBar background="bg-surface-100-900-token" padding="p0" slotTrail="space-x-0">
          <svelte:fragment slot="lead">
            <button type="button" class="btn btn-sm py-2 px-4 rounded-none">
              <IconMoreHoriz/>
            </button>
          </svelte:fragment>
          <svelte:fragment slot="trail">
            <button type="button" class="btn btn-sm py-2 px-4 hover:bg-stone-700 rounded-none">
              <IconRemove/>
            </button>
            <button type="button" class="btn btn-sm py-2 px-4 hover:bg-stone-700 rounded-none">
              <IconSelectWindowSharp/>
            </button>
            <button type="button" class="btn btn-sm py-2 px-4 hover:bg-red-700 rounded-none">
              <IconCloseRounded/>
            </button>
          </svelte:fragment>
      </AppBar>
    </div>
  </svelte:fragment>
  <div class="p-2 h-full">
    <div class="shadow-none card p-4 h-full">
      <slot />
    </div>
  </div>
</AppShell>
