<script lang="ts">
    import { goto } from '$app/navigation';
    import { ProgressBar } from '@skeletonlabs/skeleton';

    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';

    import { pack as packnamespace } from '$lib/wailsjs/go/models';
    import type { pack as packtype } from '$lib/wailsjs/go/models';

    import PackageBadge from './package-badge.svelte';
    import PackageActionButton from './package-action-button.svelte';

    const downloadHost = "download.blender.org"; // TODO: Remove this prop once the backend is ready

    export let projectId: string | undefined;
    export let pack: packtype.Package;
    export let dependencies: string[];

    let active = false;

    type BackgroundVariant = {
        variantFrom: string;
        variantTo: string;
    }

    function handleClick() {
        goto(`/packages/${pack.id}`);
    }

    function getBackgroundVariants(type: packtype.PackageType): BackgroundVariant {
        switch (type) {
            case packnamespace.PackageType.BUILD:
                return { variantFrom: 'primary', variantTo: 'secondary' };
            case packnamespace.PackageType.ADDON:
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
    on:mouseenter|stopPropagation={() => active = true}
    on:mouseleave|stopPropagation={() => active = false}
    role="button" 
    tabindex="0"
    aria-label="Interactive element"
>
    <div class="flex-shrink-0">
        <PackageActionButton 
            projectId={projectId}
            packageRef={pack.reference?.toString() || ""}
            active={active}
            variantFrom={variant.variantFrom}
            variantTo={variant.variantTo}
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
        {#if pack.state === packnamespace.PackageState.DOWNLOADING }
        <div class="py-2">
            <ProgressBar />
            <!-- <ProgressBar /> -->
        </div>
        {/if}
        <div class="flex-wrap gap-2 space-y-1 w-full">
            <PackageBadge label={downloadHost} variant="soft-success"/>
            <PackageBadge label={pack.platform?.toString()}/>
            <PackageBadge label={pack.version}/>
            <PackageBadge label={pack.author}/>
            <PackageBadge label={packnamespace.PackageType[pack.type].toLocaleLowerCase()}/>
            <PackageBadge label={packnamespace.PackageState[pack.state].toLocaleLowerCase()} />
        </div>
    </div>
</div>

