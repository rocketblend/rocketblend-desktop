<script lang="ts">
    import { PackageState, PackageType } from './types';
    import PackageListItem from './PackageListItem.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    export let packages: pack.Package[] = [];
    export let dependencies: string[] = [];

    function handleItemClick(reference: string = "") {
        console.log("click", reference);

        if (!dependencies.includes(reference)) {
            dependencies = [...dependencies, reference];
            return
        }
        
        dependencies = dependencies.filter(dep => dep !== reference);
    }

    function handleItemDownload(packageId: string = "") {
        console.log("download", packageId);
    }

    function handleItemStop(packageId: string = "") {
        console.log("stop", packageId);
    }

    function handleItemDelete(pacakgeId: string = "") {
        console.log("delete", pacakgeId);
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
            on:delete={() => handleItemDelete(pack.id?.toString())}
            on:stop={() => handleItemStop(pack.id?.toString())}
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
            platform="windows"
            downloadHost="download.blender.org"
        />
    {/each}
</div>
