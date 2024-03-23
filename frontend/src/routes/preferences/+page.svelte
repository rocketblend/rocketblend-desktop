<script lang="ts">
    import { t } from '$lib/translations/translations';
</script>

<main class="space-y-4">
    <div>
        <h2 class="h2 font-bold">{$t('preference.title')}</h2>
    </div>
    <hr>
    <div class="space-y-2">
        <h3 id="general" class="h3 font-bold">{$t('preference.category.general.title')}</h3>
        <h5 class="h5 font-bold">Projects</h5>
        <h5 class="h5 font-bold">Features</h5>
    </div>
    <div class="space-y-2">
        <h3 id="rocketblend" class="h3 font-bold">{$t('preference.category.rocketblend.title')}</h3>
    </div>
</main>