<script lang="ts">
    import type { PageData } from './$types';

    import IconDownloadFill from '~icons/ri/download-2-fill'
    import IconAddFill from '~icons/ri/add-fill';
    import IconDeleteFill from '~icons/ri/delete-bin-7-fill';
    import IconCheckFill from '~icons/ri/check-fill';
    import IconVerifiedBadgeFill from '~icons/ri/verified-badge-fill';
    import IconInfoFill from '~icons/ri/check-double-fill';

    import { pack } from '$lib/wailsjs/go/models';
    import { formatDateTime } from '$lib/components/utils';
    import { t } from '$lib/translations/translations';

    export let data: PageData;

    let active = true;
    let installed = false;

    $: installed = data.package.state === pack.PackageState.INSTALLED;
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
                <!-- <div>
                    {#if installed}
                        <button class="btn text-lg text-surface-700-200-token px-2 py-1"><IconDeleteFill/></button>
                    {/if}
                    <button class="btn text-lg variant-filled px-6 py-1"><IconAddFill/></button>
                </div> -->
            </div>
            <div class="flex flex-wrap text-sm text-surface-800-100-token gap-1">
                <div class="badge variant-ghost rounded">{data.package.reference}</div>
                <div class="badge variant-ghost rounded">{pack.PackageType[data.package.type].toLocaleLowerCase()}</div>
                <div class="badge variant-ghost rounded">{data.package.id}</div>
                <div class="badge variant-ghost rounded">{data.package.platform}</div>
                <div class="badge variant-ghost rounded">{formatDateTime(data.package.updatedAt)}</div>
                <div class="badge variant-ghost rounded">{pack.PackageState[data.package.state].toLocaleLowerCase()}</div>
            </div>
        </div>
    </div>
    <hr>

    {#if active}
        <aside class="alert variant-ghost-primary">
            <IconInfoFill class="text-2xl"/>
            <div class="alert-message ">
                <h2 class="font-bold h6">Enabled</h2>
                <p class="text-sm">Package is current enabled on the selected project</p>
            </div>
            <div class="alert-actions">
                <button class="btn btn-sm variant-glass-surface font-medium">Disable</button>
            </div>
        </aside>
    {/if}

    {#if !installed}
        <aside class="alert variant-ghost-surface">
            <IconDownloadFill class="text-2xl"/>
            <div class="alert-message ">
                <h2 class="font-bold h6">Available</h2>
                <p class="text-sm">Package is available to be downloaded.</p>
            </div>
            <div class="alert-actions">
                <button class="btn btn-sm variant-filled-surface font-medium">Download</button>
            </div>
        </aside>
    {:else}
        <aside class="alert variant-glass-success">
            <IconCheckFill class="text-2xl"/>
            <div class="alert-message ">
                <h2 class="font-bold h6">Ready</h2>
                <p class="text-sm">Package is ready to be used.</p>
            </div>
        </aside>
    {/if}

    {#if installed}
        <aside class="alert variant-ghost-surface">
            <IconDeleteFill class="text-2xl"/>
            <div class="alert-message ">
                <h2 class="font-bold h6">Delete Package</h2>
                <p class="text-sm">Package is installed and can be deleted.</p>
            </div>
            <div class="alert-actions">
                <button class="btn btn-sm variant-filled-surface font-medium">Delete</button>
            </div>
        </aside>
    {/if}
</main>