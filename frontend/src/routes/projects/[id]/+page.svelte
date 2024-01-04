<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import type { PageData } from './$types';

    import { t } from '$lib/translations/translations';
    import { getSelectedProjectStore } from '$lib/stores';
    import { debounce } from '$lib/utils';
    import { EventsOn } from '$lib/wailsjs/runtime';
    import { formatDateTime, resourcePath } from '$lib/components/utils';
    import { EVENT_DEBOUNCE, SEARCH_STORE_INSERT_CHANNEL } from '$lib/events';

	import Media from '$lib/components/core/media/Media.svelte';
	import { GetProject } from '$lib/wailsjs/go/application/Driver';
	import InlineInput from '$lib/components/core/input/InlineInput.svelte';

    import IconEditFill from '~icons/ri/edit-fill';

    const selectedProjectStore = getSelectedProjectStore();
    const refreshProjectDebounced = debounce(refreshProject, EVENT_DEBOUNCE);

    export let data: PageData;
    
    let cancelListener: () => void;

    async function refreshProject() {
        const project = (await GetProject(data.project.id?.toString())).project;
        if (!project) {
            return;
        }

        data = {...data, project};
    }

    function setSelectedProject() {
        if (data.project.id) {
            selectedProjectStore.set([data.project.id.toString()]);
        }
    }

    function getDependenciesDisplay(): string {
        const dependencies = data.project.addons || [];
        const buildDependencyCount = data.project.build ? 1 : 0;
        return t.get('home.project.tag.dependency', { number: dependencies.length + buildDependencyCount });
    }

    setSelectedProject();

    onMount(() => {
        cancelListener = EventsOn(SEARCH_STORE_INSERT_CHANNEL, (event: { id: string, indexType: string }) => {
            if (event.indexType === "project" && event.id === data.project.id?.toString()) {
                refreshProjectDebounced();
            }
        });
    });

    onDestroy(() => {
        if (cancelListener) {
            cancelListener();
        }
    });
</script>

<main class="space-y-4"> 
    <div class="flex gap-4 items-end">
        <div>
            <Media src={resourcePath(data.project.thumbnailPath)} alt="" />
        </div>
        <div class="space-y-2">
            <InlineInput bind:value={data.project.name} labelClasses="h2 font-bold items-baseline" inputClasses="input">
                <IconEditFill class="text-sm text-surface-600-300-token"/>
            </InlineInput>
            <div class="flex flex-wrap text-sm text-surface-800-100-token gap-1">
                <div class="badge variant-ghost rounded">{data.project.path}</div>
                <div class="badge variant-ghost rounded">{data.project.fileName}</div>
                <div class="badge variant-ghost rounded">{data.project.build}</div>
                {#each data.project.tags || [] as tag}
                    <div class="badge variant-ghost-primary rounded">{tag}</div>
                {/each}
                <div class="badge variant-ghost-secondary rounded">{getDependenciesDisplay()}</div>
                <div class="badge variant-ghost rounded">{formatDateTime(data.project.updatedAt)}</div>
            </div>
        </div>
    </div>
    <!-- <hr>
    <InlineInput type="textarea" placeholder="Add description..."/> -->
    <hr>
    <div class="grid grid-cols-4 gap-4">
        <Media height="80" width="full" src="{resourcePath(data.project.splashPath)}" alt="" />
      </div>
</main>