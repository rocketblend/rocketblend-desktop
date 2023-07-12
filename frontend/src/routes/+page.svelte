<script>
    import ProjectList from '$lib/components/project-list/project-list.svelte';

    import { t } from '$lib/translations/translations';
    import { Greet } from "$lib/wailsjs/go/application/Driver.js";

    let link = 'https://kit.svelte.dev';

    let resultText = "Please enter your name";
    let name = "";

    function greet() {
        Greet(name).then((result) => (resultText = result));
    }
</script>


<main class="space-y-4"> 
    <h3 class="font-bold">{$t('home.title')}</h3>
    <ProjectList/>
</main>
