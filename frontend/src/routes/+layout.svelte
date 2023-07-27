<script>
  import "../theme.postcss";
  import "@skeletonlabs/skeleton/styles/all.css";
  import "../app.postcss";

  import { AppBar, AppShell } from '@skeletonlabs/skeleton';

  import logo from '$lib/assets/images/logo.png';

  import IconCloseFill from '~icons/ri/close-fill'
  import IconMoreFill from '~icons/ri/more-fill'
  import IconSubtractFill from '~icons/ri/subtract-fill'
  import IconCheckboxMultipleBlankLine from '~icons/ri/checkbox-multiple-blank-line'
  import IconBlenderFill from '~icons/ri/blender-fill'
  import IconFolderOpenFill from '~icons/ri/folder-open-fill'
  import IconEyeFill from '~icons/ri/eye-fill'
  import IconBookmark3Fill from '~icons/ri/bookmark-3-fill'

  import { Quit, WindowMinimise, WindowToggleMaximise } from '$lib/wailsjs/runtime';
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
  <div class="h-full p-2 py-0">
    <div class="shadow-none card p-6 h-full">
      <slot />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <section class="grid grid-cols-3 gap-4 p-3 pb-3">
      <div class="flex gap-4 items-center">
        <img class="h-auto max-w-full rounded-lg" src={logo} alt="" height="64" width="64">
        <div>
          <div class="text-sm font-medium">Hello World</div>
          <div class="text-sm text-muted text-surface-300">Offical Stable Blender 3.4.1</div>
        </div>
        <button type="button" class="btn btn-lg px-1 text-secondary-300-600-token"><IconBookmark3Fill/></button>
      </div>
      <div class="min-w-max items-center justify-center flex gap-2">
        <!-- <button type="button" class="btn variant-filled text-lg px-9"><IconBlenderFill/></button> -->
      </div>
      <div class="justify-end items-center flex gap-2">
          <button type="button" class="btn text-lg text-surface-700-200-token px-2"><IconEyeFill/></button>
          <button type="button" class="btn text-lg text-surface-700-200-token px-2"><IconFolderOpenFill/></button>
          <button type="button" class="btn variant-filled text-lg px-9"><IconBlenderFill/></button>
      </div>
    </section>
  </svelte:fragment>
</AppShell>
