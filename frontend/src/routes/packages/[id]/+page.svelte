<script lang="ts">
    import type { PageData } from './$types';

    import { t } from '$lib/translations/translations';

    export let data: PageData;
</script>

<main class="flex flex-col h-full space-y-4">
    <h2 class="h2 font-bold">{data.package.name}</h2>
    <hr>
</main>