<script lang="ts">
    import { t } from '$lib/translations/translations';
    import type { PageData } from './$types';

    import ProjectList from '$lib/components/project-list/project-list.svelte';

    export let data: PageData;
</script>


<main class="space-y-4"> 
    <h3 class="font-bold">{$t('home.title')}</h3>
    <ProjectList projects={data.projects}/>
</main>
