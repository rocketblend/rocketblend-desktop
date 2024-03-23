<script lang="ts">
    import { t } from '$lib/translations/translations';
</script>

<main class="space-y-4">
    <div>
        <h2 class="h2 font-bold">{$t('preference.title')}</h2>
    </div>
    <hr>
    <div class="space-y-2">
        <h5 id="general" class="h5 font-bold">{$t('preference.category.general.title')}</h5>
    </div>
    <div class="space-y-2">
        <h5 id="rocketblend" class="h5 font-bold">{$t('preference.category.rocketblend.title')}</h5>
    </div>
    <div class="space-y-2">
        <h5 id="developer" class="h5 font-bold">{$t('preference.category.developer.title')}</h5>
    </div>
</main>