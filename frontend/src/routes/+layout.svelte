<script lang="ts">
    import "../app.postcss";

    import { AppShell, AppBar } from '@skeletonlabs/skeleton';
</script>

<AppShell slotSidebarLeft="flex flex-col overflow-y-hidden space-y-2 pl-2 w-96 h-full" slotPageContent="overflow-hidden h-full">
    <svelte:fragment slot="header">
        <div style="--wails-draggable:drag">
            <div class="h-9 p-0"></div>
        </div>
    </svelte:fragment>
    <div class="h-full p-2">
        <div class="shadow-none card p-6 h-full overflow-hidden">
            <slot />
        </div>
    </div>
</AppShell>