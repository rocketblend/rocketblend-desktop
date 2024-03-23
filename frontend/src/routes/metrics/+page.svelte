<script lang="ts">
    import { t } from '$lib/translations/translations';

    import HeatMap from '$lib/components/ui/metric/HeatMap.svelte';
</script>

<main class="space-y-4">
    <div>
        <h1 class="h2 font-bold">{$t('home.metric.title')}</h1>
    </div>

    <HeatMap />
</main>