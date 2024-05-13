<script lang="ts">
    import type { PageData } from './$types';

    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';

    import { formatDateTime } from '$lib/components/utils';
    
	import { PackageToggle, PackageActions } from './(components)';

    export let data: PageData;
</script>

<main class="flex flex-col h-full space-y-4"> 
    <div class="flex gap-4 items-end">
        <div class="space-y-2 w-full">
            <div class="flex justify-between items-center">
                <div class="flex justify-between items-baseline gap-2">
                    <h2 class="h2 font-bold capitalize">{data.package.name}</h2>
                    <h4 class="h4 text-surface-800-100-token font-medium">{data.package.tag}</h4>
                    {#if data.package.verified}
                        <IconVerifiedBadgeFill class="text-primary-500" />
                    {/if}
                </div>
            </div>
            <div class="flex flex-wrap text-sm text-surface-800-100-token gap-1">
                <div class="badge variant-ghost rounded">{data.package.reference}</div>
                <div class="badge variant-ghost rounded">{data.package.type.toLocaleLowerCase()}</div>
                <div class="badge variant-ghost rounded">{data.package.id}</div>
                <div class="badge variant-ghost rounded">{data.package.platform}</div>
                <div class="badge variant-ghost rounded">{formatDateTime(data.package.updatedAt)}</div>
                <div class="badge variant-ghost rounded">{data.package.state.toLocaleLowerCase()}</div>
            </div>
        </div>
    </div>
    <hr>

    <PackageToggle
        projectId={data.selectedProject?.project?.id.toString()}
        projectBuildRef={data.selectedProject?.project?.build}
        projectAddonRefs={data.selectedProject?.project?.addons}
        packageRef={data.package.reference || ''}
        packageState={data.package.state}
    />
    <PackageActions
        state={data.package.state}
        packageId={data.package.id.toString()}
        progress={data.package.progress}
    />
</main>