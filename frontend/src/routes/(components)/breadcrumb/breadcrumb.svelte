<script lang="ts">
    import { page } from '$app/stores';

    let crumbs: Array<{ label: string, href: string }> = [];

    $: {
        const tokens = $page.url.pathname.split('/').filter((t) => t !== '');

        let tokenPath = '';
        crumbs = tokens.map((t, index) => {
            tokenPath += '/' + t;

            const formattedToken = t.charAt(0).toUpperCase() + t.slice(1);
            const isCurrentPage = index === tokens.length - 1;
            const label = isCurrentPage && $page.data.label ? $page.data.label : formattedToken;

            return {
                label: label,
                href: tokenPath
            };
        });

        crumbs.unshift({ label: 'Home', href: '/' });
    }
</script>

{#if $page.data.showBreadcrumb}
    <ol class="breadcrumb text-sm text-surface-800-100-token truncate">
        {#each crumbs as crumb, i}
            {#if i < crumbs.length - 1}
                <li class="crumb"><a href={crumb.href}>{crumb.label}</a></li>
                <li class="crumb-separator" aria-hidden>&rsaquo;</li>
            {:else}
                <li class="crumb font-medium">{crumb.label}</li>
            {/if}
        {/each}
    </ol>
{/if}
