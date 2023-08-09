<script lang="ts">
    import { goto } from '$app/navigation';

    import IconBlenderFill from '~icons/ri/blender-fill'
    import IconFolderOpenFill from '~icons/ri/folder-open-fill'
    import IconEyeFill from '~icons/ri/eye-fill'
    import IconBookmark3Fill from '~icons/ri/bookmark-3-fill'

    import { FxReveal as Img } from '@zerodevx/svelte-img'

    import logo from '$lib/assets/images/logo.png?run';

    import { selectedProject } from '$lib/store';

    function handleViewProject(): void {
        goto(`/projects/${$selectedProject?.id}`);
    }
</script>

<section class="grid grid-cols-3 gap-4 p-3 pb-3">
    {#if $selectedProject}
        <div class="flex gap-4 items-center">
            <Img class="h-auto max-w-full rounded-lg w-16 image-fade-in" src={logo} alt="" />
            <div>
                <div class="text-sm font-medium">{$selectedProject.name}</div>
                <div class="text-sm text-surface-300">{$selectedProject.fileName}</div>
            </div>
            <button type="button" class="btn btn-lg px-1 text-secondary-300-600-token"><IconBookmark3Fill/></button>
        </div>
        <div class="min-w-max items-center justify-center flex gap-2"></div>
        <div class="justify-end items-center flex gap-2">
            <button type="button" class="btn text-lg text-surface-700-200-token px-2" on:click={handleViewProject}><IconEyeFill/></button>
            <button type="button" class="btn text-lg text-surface-700-200-token px-2"><IconFolderOpenFill/></button>
            <button type="button" class="btn variant-filled text-lg px-9"><IconBlenderFill/></button>
        </div>
    {:else}
        <div class="flex gap-4 items-center">
            <div class="rounded-lg placeholder animate-pulse h-16 w-16"></div>
            <div>
                <div class="placeholder animate-pulse h-4 w-24 mb-2"></div>
                <div class="placeholder animate-pulse h-4 w-32"></div>
            </div>
            <!-- <div class="placeholder animate-pulse px-1 h-6 w-5"></div> -->
        </div>
        <div class="min-w-max items-center justify-center flex gap-2"></div>
        <div class="justify-end items-center flex gap-4">
            <!-- <div class="placeholder animate-pulse w-6 h-6"></div>
            <div class="placeholder animate-pulse w-6 h-6"></div> -->
            <div class="placeholder animate-pulse w-24 h-10"></div>
        </div>
    {/if}
</section>