<script lang="ts">
    import { ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
    import type { pack } from '$lib/wailsjs/go/models';
  
    export let packages: pack.Package[] = [];
    export let selectedPackageIds: string[] = [];
</script>

<ListBox class="flex-auto" multiple>
    {#each packages as pack}
        <ListBoxItem bind:group={selectedPackageIds} name="packages" value={pack.id} active="variant-glass-primary" hover="hover:variant-filled-surface" rounded="rounded" class="truncate overflow-hidden">
            <div class="flex gap-2">
                <!-- <div class="flex-shrink-0 inline-flex items-center">
                    {#if pack.type == 'build'}
                        <IconCheckboxBlankCircleLine />
                    {:else}
                        <IconCheckboxBlankLine />
                    {/if}
                </div> -->
                <div class="flex-grow w-2">
                    <div class="text-sm font-medium">{pack.name}</div>
                    <div class="text-sm">{pack.reference}</div>
                </div>
            </div>
        </ListBoxItem>
    {/each}
</ListBox>

