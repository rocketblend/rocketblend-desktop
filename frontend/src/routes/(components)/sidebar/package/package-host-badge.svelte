<script lang="ts">
    import PackageBadge from './package-badge.svelte';
  
    // TODO: Remove this prop once the backend is ready
    const trustedHosts = ["download.blender.org", "github.com", "gitlab.com", "bitbucket.org"];

    export let uri: string;
    let hostName: string = '';
    let safe: boolean;
  
    $: {
        try {
            hostName = new URL(uri).hostname;
            safe = trustedHosts.includes(hostName);
        } catch (error) {
            hostName = '';
            safe = false;
        }
    }
    $: variant = safe ? "soft-success" : undefined;
  </script>
  
  {#if hostName}
    <PackageBadge label={hostName} variant={variant} />
  {/if}
  