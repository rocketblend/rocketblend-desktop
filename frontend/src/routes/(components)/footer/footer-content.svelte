<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    import IconBlenderFill from '~icons/ri/blender-fill';
    import IconFolderOpenFill from '~icons/ri/folder-open-fill';
    import IconEyeFill from '~icons/ri/eye-fill';
    
    import { Media } from '$lib/components/ui/gallery';

    const dispatch = createEventDispatcher();

    export let name: string = '';
    export let fileName: string = '';
    export let imagePath: string = '';
    export let isLoading: boolean = true;
    export let disabled: boolean = false;

    function handleViewProject() {
        dispatch('viewProject');
    }

    function handleRunProject() {
        dispatch('runProject');
    }

    function handleExploreProject() {
        dispatch('exploreProject');
    }
</script>

<section class="grid grid-cols-3 gap-4 p-3 pb-3">
    {#if !isLoading}
        <div class="flex gap-4 items-center max-h-16">
            <Media src={imagePath} height={16} width={16} class="cursor-default" rounded>
                <span slot="not-found" class="text-sm">?</span>
            </Media>
            <div>
                <div class="text-sm font-medium">{name}</div>
                <div class="text-sm text-surface-300">{fileName}</div>
            </div>
            <!-- <button type="button" class="btn btn-lg px-1 text-secondary-300-600-token"><IconBookmark3Fill/></button> -->
        </div>
        <div class="min-w-max items-center justify-center flex gap-2"></div>
        <div class="justify-end items-center flex gap-2">
            <button type="button" class="btn text-lg text-surface-700-200-token px-2" on:click={handleViewProject}><IconEyeFill/></button>
            <button type="button" class="btn text-lg text-surface-700-200-token px-2" on:click={handleExploreProject}><IconFolderOpenFill /></button>
            <button type="button" class="btn variant-filled text-lg px-9" on:click={handleRunProject} disabled={disabled}><IconBlenderFill/></button>
        </div>
    {:else}
        <div class="flex gap-4 items-center">
            <div class="rounded-lg placeholder animate-pulse h-16 w-16"></div>
            <div>
                <div class="placeholder animate-pulse h-4 w-24 mb-2"></div>
                <div class="placeholder animate-pulse h-4 w-32"></div>
            </div>
            <!-- <div class="placeholder animate-pulse px-1 h-6 w-5"></div> -->
        </div>
        <div class="min-w-max items-center justify-center flex gap-2"></div>
        <div class="justify-end items-center flex gap-4">
            <!-- <div class="placeholder animate-pulse w-6 h-6"></div>
            <div class="placeholder animate-pulse w-6 h-6"></div> -->
            <div class="placeholder animate-pulse w-24 h-10"></div>
        </div>
    {/if}
</section>