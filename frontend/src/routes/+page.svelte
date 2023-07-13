<script>
    import ProjectList from '$lib/components/project-list/project-list.svelte';

    import { t } from '$lib/translations/translations';
    import { Greet } from "$lib/wailsjs/go/application/Driver.js";
</script>


<main class="space-y-4"> 
    <h3 class="font-bold">{$t('home.title')}</h3>
    <ProjectList/>
</main>
