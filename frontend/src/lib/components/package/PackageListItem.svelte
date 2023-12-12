<script lang="ts">
    import { PackageState, PackageType } from './types';
    import ProgressBar from '$lib/components/core/progress/ProgressBar.svelte';
    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';
	import PackageBadge from './PackageBadge.svelte';
	import PackageActionButton from './PackageActionButton.svelte';

    export let name = "";
    export let tag = "";
    export let version = "";
    export let author = "";
    export let platform = "";
    export let reference = "";
    export let type: PackageType = PackageType.Unknown;
    export let state: PackageState = PackageState.Available;
    export let verified = false;
    export let progress = 0;
    export let downloadHost = "";
    export let selected = false;

    let active = false;

    function getBackgroundVariants(type: PackageType) {
        switch (type) {
            case PackageType.Build:
                return { variantFrom: 'primary', variantTo: 'secondary' };
            case PackageType.Addon:
                return { variantFrom: 'tertiary', variantTo: 'primary' };
            default:
                return { variantFrom: 'secondary', variantTo: 'tertiary' };
        }
    }

    const variant = getBackgroundVariants(type);

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
    aria-label="Interactive element"
>
    <div class="flex-shrink-0">
        <PackageActionButton 
            on:delete
            on:download
            on:stop
            state={state}
            isOpen={active}
            variantFrom={variant.variantFrom}
            variantTo={variant.variantTo}
        />
    </div>
    <div
        class="flex-col gap-2 overflow-hidden"
    >
        <div class="inline-flex items-center gap-2 w-full">
            <span class="font-medium truncate">{name}</span>
            <span class="text-sm truncate">{tag}</span>
            {#if verified}
                <IconVerifiedBadgeFill class="text-sm text-primary-500" />
            {/if}
        </div>
        {#if state === PackageState.Downloading }
            <ProgressBar value={progress} rounded={true} />
        {/if}
        <div class="text-sm text-surface-800-100-token truncate">{reference}</div>
        <div class="flex-wrap gap-2 space-y-1 w-full">
            <PackageBadge label={downloadHost} variant="soft-success"/>
            <PackageBadge label={platform}/>
            <PackageBadge label={version}/>
            <PackageBadge label={author}/>
            <PackageBadge label={PackageType[type].toLocaleLowerCase()}/>
            <PackageBadge label={PackageState[state].toLocaleLowerCase()}/>
        </div>
    </div>
</div>

