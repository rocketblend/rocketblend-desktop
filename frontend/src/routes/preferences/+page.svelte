<script lang="ts">
    import type { PageData } from "./$types";

    import { t } from '$lib/translations/translations';

	import { GeneralForm } from './(components)/forms';

    export let data: PageData;
</script>

<main class="space-y-4">
    <div>
        <h2 class="h2 font-bold">{$t('preference.title')}</h2>
    </div>
    <hr>
    <div class="space-y-4">
        <h3 id="general" class="h3 font-bold">{$t('preference.category.general.title')}</h3>
        <GeneralForm data={data.form}/>
    </div>
    <div class="space-y-2">
        <h3 id="rocketblend" class="h3 font-bold">{$t('preference.category.rocketblend.title')}</h3>
    </div>
</main>