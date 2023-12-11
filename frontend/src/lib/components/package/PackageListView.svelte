<script lang="ts">
    import PackageListItem from './PackageListItem.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    export let packages: pack.Package[] = [];
    export let dependencies: string[] = [];

    function handleClick(reference: string = "") {
        console.log("click", reference);

        if (!dependencies.includes(reference)) {
            dependencies = [...dependencies, reference];
            return
        }
        
        dependencies = dependencies.filter(dep => dep !== reference);
    }

    function handleDownload(packageId: string = "") {
        console.log("download", packageId);
    }

    function handleStop(packageId: string = "") {
        console.log("stop", packageId);
    }

    function handleDelete(pacakgeId: string = "") {
        console.log("delete", pacakgeId);
    }
</script>

<div class="flex-auto space-y-1 rounded-token">
    {#each packages as pack}
        <PackageListItem
            on:click={() => handleClick(pack.reference?.toString())}
            on:download={() => handleDownload(pack.id?.toString())}
            on:delete={() => handleDelete(pack.id?.toString())}
            on:stop={() => handleStop(pack.id?.toString())}
            selected={dependencies.includes(pack.reference?.toString() || "")}
            name={pack.name}
            tag={pack.tag}
            version={pack.version}
            author={pack.author}
            reference={pack.reference}
            progress={pack.installationPath ? 100 : 0}
            type={pack.type}
            verified={pack.verified}
            platform="windows"
            downloadHost="download.blender.org"
            state={pack.installationPath ? "installed" : "available"}
        />
    {/each}
</div>
