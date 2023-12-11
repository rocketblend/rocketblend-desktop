<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { ProgressBar } from '@skeletonlabs/skeleton';

    import IconDownload2Fill from '~icons/ri/file-download-fill';
    import IconStopFill from '~icons/ri/stop-mini-fill';
    import IconSubtractFill from '~icons/ri/subtract-fill';
    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';

    import IconAddFill from '~icons/ri/add-fill';

    import IconMoreFill from '~icons/ri/delete-bin-7-fill';

    export let name: string = "";
    export let tag: string = "";
    export let version: string = "";
    export let author: string = "";
    export let platform: string = "";
    export let reference: string = "";
    export let type: number | undefined;
    export let verified: boolean = false;
    export let progress: number = 0;
    export let downloadHost: string = "";
    export let state: string = "";

    export let selected: boolean = false;

    let selectedClass: string;
    let active = false;

    const dispatch = createEventDispatcher();

    const bageBackground: Record<string, string> = {
        build: "variant-gradient-primary-secondary",
        addon: "variant-gradient-tertiary-primary",
        unknown: "variant-gradient-secondary-tertiary",
    }

    function HandleClick() {
        if (progress === 0) {
            dispatch('download');
        } else if (progress > 0 && progress < 100) {
            dispatch('stop');
        } else if (progress === 100) {
            dispatch('delete');
        }
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            if (progress === 0) {
                dispatch('download');
            } else if (progress > 0 && progress < 100) {
                dispatch('stop');
            } else if (progress === 100) {
                dispatch('delete');
            }
        }
    }

    $: selectedClass = selected ? "variant-ghost-primary" : "hover:variant-filled-surface";
</script>

<div class="flex gap-2 rounded p-2 {selectedClass}"
    on:click|stopPropagation
    on:dblclick|stopPropagation
    on:keydown|stopPropagation
    on:mouseenter|stopPropagation={() => active = true}
    on:mouseleave|stopPropagation={() => active = false}
    role="button" 
    tabindex="0"
>
    <div class="flex-shrink-0">
        <div 
            class="flex items-center h-full bg-gradient-to-br {bageBackground[type || 'unknown']} rounded p-1 text-token"
            on:click|stopPropagation={HandleClick}
            on:keydown|stopPropagation={handleKeyDown}
            role="button" 
            tabindex="0"
        >
            {#if active}
                <div>
                    {#if progress == 0}
                        <IconDownload2Fill />
                    {:else if progress < 100}
                        <IconStopFill />
                    {:else}
                        <!-- {#if selected}
                            <IconSubtractFill />
                        {:else}
                            <IconAddFill />
                        {/if} -->
                        <IconMoreFill />
                    {/if}
                </div>
            {/if}
        </div>
    </div>
    <div
        class="flex-col gap-2 overflow-hidden"
    >
        <!-- Render package details -->
        <div class="inline-flex items-center gap-2 w-full">
            <span class="font-medium truncate">{name}</span>
            <span class="text-sm truncate">{tag}</span>
            {#if verified}
                <IconVerifiedBadgeFill class="text-sm text-primary-500" />
            {/if}
        </div>
        {#if progress && progress != 100 }
            <div class="flex items-center gap-2">
                <ProgressBar value={progress} rounded="rounded"/>
                <div class="text-surface-800-100-token text-sm">{progress}%</div>
            </div>
        {/if}
        <div class="text-sm text-surface-800-100-token truncate">{reference}</div>
        <div class="flex-wrap gap-2 space-y-1 w-full">
            <div class="badge variant-soft-success rounded">{downloadHost}</div>
            <div class="badge variant-ghost">{platform}</div>
            <div class="badge variant-ghost rounded">{version}</div>
            <div class="badge variant-ghost rounded">{author}</div>
            <div class="badge variant-ghost rounded">{type}</div>
            <div class="badge variant-ghost rounded">{state}</div>
        </div>
    </div>
</div>

