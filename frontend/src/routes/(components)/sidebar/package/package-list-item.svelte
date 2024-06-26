<script lang="ts">
    import { goto } from '$app/navigation';


    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';

    import { enums } from '$lib/wailsjs/go/models';
    import type { types } from '$lib/wailsjs/go/models';

    import PackageBadge from './package-badge.svelte';
    import PackageHostBadge from './package-host-badge.svelte';
    import PackageActionButton from './package-action-button.svelte';
	import ProgressBar from '$lib/components/ui/progress/progress-bar.svelte';

    export let projectId: string | undefined;
    export let pack: types.Package;
    export let dependencies: string[];

    let hovered = false;

    type BackgroundVariant = {
        variantFrom: string;
        variantTo: string;
    }

    function handleClick() {
        goto(`/packages/${pack.id}`);
    }

    function getBackgroundVariants(type: enums.PackageType): BackgroundVariant {
        switch (type) {
            case enums.PackageType.BUILD:
                return { variantFrom: 'primary', variantTo: 'secondary' };
            case enums.PackageType.ADDON:
                return { variantFrom: 'tertiary', variantTo: 'primary' };
            default:
                return { variantFrom: 'secondary', variantTo: 'tertiary' };
        }
    }

    $: selected = dependencies.includes(pack.reference?.toString() || "");
    $: variant = getBackgroundVariants(pack.type);
    $: selectedClass = selected ? "variant-ghost-primary" : "hover:variant-filled-surface";
</script>

<div class="flex gap-2 rounded p-2 {selectedClass}"
    on:click|stopPropagation={handleClick}
    on:dblclick|stopPropagation
    on:keydown|stopPropagation
    on:mouseenter|stopPropagation={() => hovered = true}
    on:mouseleave|stopPropagation={() => hovered = false}
    role="button" 
    tabindex="0"
    aria-label="Interactive element"
>
    <div class="flex-shrink-0">
        <PackageActionButton
            projectId={projectId}
            packageRef={pack.reference?.toString() || ""}
            hovered={hovered}
            variantFrom={variant.variantFrom}
            variantTo={variant.variantTo}
            assigned={selected}
        />
    </div>
    <div
        class="flex-col gap-2 overflow-hidden"
    >
        <div class="inline-flex items-center gap-2 w-full">
            <span class="font-medium truncate">{pack.name}</span>
            <span class="text-sm truncate">{pack.tag}</span>
            {#if !!pack.verified}
                <IconVerifiedBadgeFill class="text-sm text-primary-500" />
            {/if}
        </div>
        <div class="text-sm text-surface-800-100-token truncate">{pack.reference}</div>
        {#if pack.state === enums.PackageState.DOWNLOADING }
        <div class="py-2">
            <ProgressBar value={pack.progress?.currentBytes} max={pack.progress?.totalBytes} />
            <!-- <ProgressBar /> -->
        </div>
        {/if}
        <div class="flex-wrap gap-2 space-y-1 w-full">
            <PackageHostBadge uri={pack.uri}/>
            <PackageBadge label={pack.platform.toString()}/>
            {#if pack.version }
                <PackageBadge label={pack.version}/>
            {/if}
            <PackageBadge label={pack.author}/>
            <PackageBadge label={pack.type.toLocaleLowerCase()}/>
            <PackageBadge label={pack.state.toLocaleLowerCase()} />
        </div>
    </div>
</div>

