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

<AppShell regionPage="p-2">
  <svelte:fragment slot="header">
    <div style="--wails-draggable:drag">
      <AppBar background="bg-surface-100-900-token" padding="px-4 pt-2 pb-0" slotTrail="space-x-7">
          <svelte:fragment slot="lead"><IconMoreHoriz class="text-2xl"/></svelte:fragment>
          <svelte:fragment slot="trail"><IconRemove/><IconSelectWindowSharp/><IconCloseRounded/></svelte:fragment>
      </AppBar>
    </div>
  </svelte:fragment>
  <div class="shadow-none card p-4 h-full">
    <slot />
  </div>
</AppShell>
