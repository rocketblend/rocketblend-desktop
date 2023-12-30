<script lang="ts">
    import ProgressBar from '$lib/components/core/progress/ProgressBar.svelte';
    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';
	import PackageBadge from './PackageBadge.svelte';
	import PackageActionButton from './PackageActionButton.svelte';
	import { pack } from '$lib/wailsjs/go/models';

    const platform = "windows"; // TODO: Remove this prop once the backend is ready
    const downloadHost = "download.blender.org"; // TODO: Remove this prop once the backend is ready

    export let item: pack.Package;
    export let selected = false;

    let active = false;

    function getBackgroundVariants(type: pack.PackageType) {
        switch (type) {
            case pack.PackageType.BUILD:
                return { variantFrom: 'primary', variantTo: 'secondary' };
            case pack.PackageType.ADDON:
                return { variantFrom: 'tertiary', variantTo: 'primary' };
            default:
                return { variantFrom: 'secondary', variantTo: 'tertiary' };
        }
    }

    const variant = getBackgroundVariants(item.type || 0);

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
            state={item.state}
            isOpen={active}
            variantFrom={variant.variantFrom}
            variantTo={variant.variantTo}
        />
    </div>
    <div
        class="flex-col gap-2 overflow-hidden"
    >
        <div class="inline-flex items-center gap-2 w-full">
            <span class="font-medium truncate">{item.name}</span>
            <span class="text-sm truncate">{item.tag}</span>
            {#if item.verified}
                <IconVerifiedBadgeFill class="text-sm text-primary-500" />
            {/if}
        </div>
        {#if item.state === pack.PackageState.DOWNLOADING }
            <ProgressBar rounded={true} />
        {/if}
        <div class="text-sm text-surface-800-100-token truncate">{item.reference}</div>
        <div class="flex-wrap gap-2 space-y-1 w-full">
            <PackageBadge label={downloadHost} variant="soft-success"/>
            <PackageBadge label={item.platform?.toString()}/>
            <PackageBadge label={item.version}/>
            <PackageBadge label={item.author}/>
            <PackageBadge label={pack.PackageType[item.type || 0].toLocaleLowerCase()}/>
            <PackageBadge label={pack.PackageState[item.state || 0].toLocaleLowerCase()}/>
        </div>
    </div>
</div>

