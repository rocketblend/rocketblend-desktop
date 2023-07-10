<script>
  import "../theme.postcss";
  import "@skeletonlabs/skeleton/styles/all.css";
  import "../app.postcss";

  import { AppBar, AppShell } from '@skeletonlabs/skeleton';
</script>

<AppShell regionPage="p-2">
  <svelte:fragment slot="header">
    <div style="--wails-draggable:drag">
      <AppBar background="bg-surface-100-900-token">
          <svelte:fragment slot="lead">(...)</svelte:fragment>
          <svelte:fragment slot="trail">_ [] X</svelte:fragment>
      </AppBar>
    </div>
  </svelte:fragment>
  <div class="shadow-none card p-4 h-full">
    <slot />
  </div>
</AppShell>
