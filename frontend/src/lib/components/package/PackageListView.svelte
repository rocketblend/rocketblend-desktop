<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { PackageState, PackageType } from './types';
    import PackageListItem from './PackageListItem.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    const dispatch = createEventDispatcher();

    export let packages: pack.Package[] = [];
    export let dependencies: string[] = [];
    export let platform: string = "windows";

    function handleItemClick(reference: string = "") {
        console.log("click", reference);

        if (!dependencies.includes(reference)) {
            dependencies = [...dependencies, reference];
            return
        }
        
        dependencies = dependencies.filter(dep => dep !== reference);
    }

    function handleItemDownload(packageId: string = "") {
        dispatch('download', { packageId });
    }

    function handleItemCancel(packageId: string = "") {
        dispatch('cancel', { packageId });
    }

    function handleItemDelete(packageId: string = "") {
        dispatch('delete', { packageId });
    }

    function getPackageType(typeStr: string = ""): PackageType {
        switch (typeStr) {
            case "build":
                return PackageType.Build;
            case "addon":
                return PackageType.Addon;
            default:
                return PackageType.Unknown;
        }
    }

    function getPackageState(installationPath: string = ""): PackageState {
        if (installationPath) {
            return PackageState.Installed;
        }

        return PackageState.Available;
    }
</script>

<div class="flex-auto space-y-1 rounded-token">
    {#each packages as pack}
        <PackageListItem
            on:click={() => handleItemClick(pack.reference?.toString())}
            on:download={() => handleItemDownload(pack.id?.toString())}
            on:cancel={() => handleItemCancel(pack.id?.toString())}
            on:delete={() => handleItemDelete(pack.id?.toString())}
            selected={dependencies.includes(pack.reference?.toString() || "")}
            name={pack.name}
            type={getPackageType(pack.type?.toString())}
            state={getPackageState(pack.installationPath?.toString())}
            tag={pack.tag}
            version={pack.version}
            author={pack.author}
            reference={pack.reference}
            progress={pack.installationPath ? 100 : 0}
            verified={pack.verified}
            platform={platform}
            downloadHost="download.blender.org"
        />
    {/each}
</div>
