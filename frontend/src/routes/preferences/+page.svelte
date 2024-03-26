<script lang="ts">
    import { invalidateAll } from '$app/navigation';
    import type { PageData } from "./$types";

    import { SlideToggle } from "@skeletonlabs/skeleton";

    import { t } from '$lib/translations/translations';
    import { application } from "$lib/wailsjs/go/models";
    import { OpenDirectoryDialog, OpenExplorer, UpdatePreferences } from '$lib/wailsjs/go/application/Driver';

    export let data: PageData;

    async function savePreferences() {
        try {
            const opts = new application.UpdatePreferencesOpts({
                watchPath: data.preferences.watchPath,
                feature: data.preferences.feature
            });

            await UpdatePreferences(opts);
            await refreshData();
        } catch (err) {
            console.error(err);
        }
    }

    async function refreshData() {
        await invalidateAll();
    }

    async function handleChangeProjectWatchDirectory() {
        const opts = new application.OpenDialogOptions({
            title: "Select project watch directory",
        });

        const result = await OpenDirectoryDialog(opts);
        if (!result) {
            return;
        }

        data.preferences.watchPath = result;
        await savePreferences();
    }

    async function handleExplore(path: string) {
        if (!path) return;

        const opts = new application.OpenExplorerOptions({ path });
        await OpenExplorer(opts);
    }
</script>

<main class="h-full overflow-auto space-y-4">
    <div>
        <h2 class="h2 font-bold">{$t('preference.title')}</h2>
    </div>
    <hr>
    <div class="space-y-8">
        <div class="space-y-3">
            <h5 class="h6 font-bold">General</h5>
            <div class="flex justify-between items-center gap-6">
                <div class="text-sm text-left">
                    <span class="font-medium">Project watch directory</span><br>
                    <span class="text-surface-200 ">{data.preferences.watchPath}</span>
                </div>
                <button class="btn variant-filled-surface text-sm font-medium" on:click={handleChangeProjectWatchDirectory}>
                    Change location
                </button>
            </div>

            <div class="flex justify-between items-center gap-6">
                <div class="text-sm text-left">
                    <span class="font-medium">Configuration file</span><br>
                    <span class="text-surface-200 ">{data.details.applicationConfigPath}</span>
                </div>
                <button class="btn variant-filled-surface text-sm font-medium" on:click={() => { handleExplore(data.details.applicationConfigPath);}}>
                    View location
                </button>
            </div>
        </div>

        <div class="space-y-3">
            <h5 class="h6 font-bold">RocketBlend</h5>
            <div class="flex justify-between items-center gap-6">
                <div class="text-sm text-left">
                    <span class="font-medium">Configuration file</span><br>
                    <span class="text-surface-200 ">{data.details.rocketblendConfigPath}</span>
                </div>
                <button class="btn variant-filled-surface text-sm font-medium" on:click={() => { handleExplore(data.details.rocketblendConfigPath);}}>
                    View location
                </button>
            </div>
        </div>

        <div class="space-y-3">
            <div>
                <h5 class="h6 font-bold">Experimental</h5>
                <p class="text-sm text-surface-200">
                    Enable experimental features to try out new functionalities before they're fully released.
                    <span class="font-medium">Please note, these features might be unstable or incomplete.</span>
                </p>
            </div>
            <div class="flex justify-between items-center gap-6">
                <div class="text-sm text-left">
                    Addons - <span class="text-surface-200">Configure and install blender addons for projects.</span>
                </div>
                <SlideToggle
                    name="addon"
                    size="sm"
                    active="bg-secondary-500"
                    bind:checked={data.preferences.feature.addon}
                    on:change={savePreferences}
                />
            </div>
            <div class="flex justify-between items-center gap-6">
                <div class="text-sm text-left">
                    Developer mode - <span class="text-surface-200">Access advanced settings and features for development purposes.</span>
                </div>
                <SlideToggle
                    name="developer"
                    size="sm"
                    active="bg-secondary-500"
                    bind:checked={data.preferences.feature.developer}
                    on:change={savePreferences}
                />
            </div>
        </div>
    </div>
</main>