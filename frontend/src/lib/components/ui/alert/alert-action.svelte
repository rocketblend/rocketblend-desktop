<script lang="ts">
    export let weight: "medium" | "bold" = "medium";
    export let variant: "filled-surface" | "glass-surface" = "filled-surface";
    export let disabled: boolean = false;
    export let href: string = "";

    $: disable = disabled ? "disabled" : "";
</script>


{#if href}
    <a class="btn btn-sm font-{weight} variant-{variant} {disable}" href={href}>
        <slot>Action</slot>
    </a>
{:else}
    <button class="btn btn-sm font-{weight} variant-{variant}" disabled={disabled} on:click>
        <slot>Action</slot>
    </button>
{/if}