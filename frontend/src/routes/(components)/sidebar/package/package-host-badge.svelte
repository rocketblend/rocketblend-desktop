<script lang="ts">
    import PackageBadge from './package-badge.svelte';
  
    // TODO: Remove this once the backend is ready
    const trustedHosts = ["download.blender.org", "github.com", "gitlab.com", "bitbucket.org"];

    export let uri: string;
    let hostName: string = '';
    let trusted: boolean;
  
    $: {
        try {
            hostName = new URL(uri).hostname;
            trusted = trustedHosts.includes(hostName);
        } catch (error) {
            hostName = '';
            trusted = false;
        }
    }
    $: variant = trusted ? "soft-success" : undefined;
  </script>
  
  {#if hostName}
    <PackageBadge label={hostName} variant={variant} />
  {/if}
  