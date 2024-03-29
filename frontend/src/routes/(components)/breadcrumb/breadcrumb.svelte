<script lang="ts">
    import { page } from '$app/stores';

    let crumbs: Array<{ label: string, href: string }> = [];

    $: {
        // Remove zero-length tokens.
        const tokens = $page.url.pathname.split('/').filter((t) => t !== '');

        // Create { label, href } pairs for each token.
        let tokenPath = '';
        crumbs = tokens.map((t) => {
        tokenPath += '/' + t;
        t = t.charAt(0).toUpperCase() + t.slice(1);
            return {
                label: $page.data.label || t,
                href: tokenPath
            };
        });

        // Add a way to get home too.
        crumbs.unshift({ label: 'Home', href: '/' });
    }
</script>

{#if $page.data.showBreadcrumb}
    <ol class="breadcrumb text-sm text-surface-800-100-token truncate">
        {#each crumbs as crumb, i}
            <!-- If crumb index is less than the breadcrumb length minus 1 -->
            {#if i < crumbs.length - 1}
                <li class="crumb"><a href={crumb.href}>{crumb.label}</a></li>
                <li class="crumb-separator" aria-hidden>&rsaquo;</li>
            {:else}
                <li class="crumb font-medium">{crumb.label}</li>
            {/if}
        {/each}
    </ol>
{/if}
