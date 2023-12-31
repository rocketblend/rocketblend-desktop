<script lang="ts">
    import ProgressBar from '$lib/components/core/progress/ProgressBar.svelte';
    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';
    import PackageBadge from './PackageBadge.svelte';
    import PackageActionButton from './PackageActionButton.svelte';
    import { pack } from '$lib/wailsjs/go/models';

    const downloadHost = "download.blender.org"; // TODO: Remove this prop once the backend is ready

    export let name: string;
    export let tag: string;
    export let verified: boolean;
    export let reference: string;
    export let platform: string;
    export let version: string;
    export let author: string;
    export let type: pack.PackageType = pack.PackageType.UNKNOWN;
    export let state: pack.PackageState = pack.PackageState.ERROR;
    export let selected = false;

    let active = false;

    type BackgroundVariant = {
        variantFrom: string;
        variantTo: string;
    }

    function getBackgroundVariants(type: pack.PackageType): BackgroundVariant {
        switch (type) {
            case pack.PackageType.BUILD:
                return { variantFrom: 'primary', variantTo: 'secondary' };
            case pack.PackageType.ADDON:
                return { variantFrom: 'tertiary', variantTo: 'primary' };
            default:
                return { variantFrom: 'secondary', variantTo: 'tertiary' };
        }
    }

    $: variant = getBackgroundVariants(type);
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
            on:download
            on:cancel
            on:delete
            state={state}
            open={active}
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
        {#if state === pack.PackageState.DOWNLOADING }
            <ProgressBar rounded={true} />
        {/if}
        <div class="text-sm text-surface-800-100-token truncate">{reference}</div>
        <div class="flex-wrap gap-2 space-y-1 w-full">
            <PackageBadge label={downloadHost} variant="soft-success"/>
            <PackageBadge label={platform?.toString()}/>
            <PackageBadge label={version}/>
            <PackageBadge label={author}/>
            <PackageBadge label={pack.PackageType[type].toLocaleLowerCase()}/>
            <PackageBadge label={pack.PackageState[state].toLocaleLowerCase()}/>
        </div>
    </div>
</div>

