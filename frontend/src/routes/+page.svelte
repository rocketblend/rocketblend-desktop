<script lang="ts">
    import { t } from '$lib/translations/translations';
    import ProjectList from '$lib/components/project-list/project-list.svelte';
</script>


<main class="space-y-4"> 
    <h3 class="font-bold">{$t('home.title')}</h3>
    <ProjectList />
</main>
