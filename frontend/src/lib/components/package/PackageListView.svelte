<script lang="ts">
    import { ListBox, ListBoxItem } from '@skeletonlabs/skeleton';
    import PackageListItem from './PackageListItem.svelte';
    import type { pack } from '$lib/wailsjs/go/models';

    export let packages: pack.Package[] = [];
    export let selectedPackageIds: string[] = [];

    function handleDownload(id: string = "") {
        console.log("download", id);
    }

    function handleStop(id: string = "") {
        console.log("stop", id);
    }

    function handleAdd(id: string = "") {
        console.log("add", id);
    }

    function handleRemove(id: string = "") {
        console.log("remove", id);
    }

</script>

<ListBox class="flex-auto" regionDefault="w-full" multiple>
    {#each packages as pack}
        <ListBoxItem bind:group={selectedPackageIds} name="packages" value={pack.id} active="variant-glass-primary" hover="hover:variant-filled-surface" rounded="rounded">
            <PackageListItem
                on:download={() => handleDownload(pack.id?.toString())}
                on:stop={() => handleStop(pack.id?.toString())}
                on:add={() => handleAdd(pack.id?.toString())}
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
        </ListBoxItem>
    {/each}
</ListBox>
