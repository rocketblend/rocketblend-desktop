<script lang="ts">
    import PackageListItem from './PackageListItem.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    export let packages: pack.Package[] = [];
    export let dependencies: string[] = [];

    function handleClick(packageId: string = "") {
        console.log("click", packageId);
    }

    function handleDownload(packageId: string = "") {
        console.log("download", packageId);
    }

    function handleStop(packageId: string = "") {
        console.log("stop", packageId);
    }

    function handleAdd(reference: string = "") {
        if (!dependencies.includes(reference)) {
            dependencies = [...dependencies, reference];
        }

        console.log("add", reference);
    }

    function handleRemove(reference: string = "") {
        dependencies = dependencies.filter(dep => dep !== reference);
        console.log("remove", reference);
    }
</script>

<div class="flex-auto space-y-1 rounded-token">
    {#each packages as pack}
        <PackageListItem
            on:click={() => handleClick(pack.id?.toString())}
            on:download={() => handleDownload(pack.id?.toString())}
            on:stop={() => handleStop(pack.id?.toString())}
            on:add={() => handleAdd(pack.reference?.toString())}
            on:remove={() => handleRemove(pack.reference?.toString())}
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
