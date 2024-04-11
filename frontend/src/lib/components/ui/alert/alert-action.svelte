<script lang="ts">
    export let text: string;
    export let weight: "medium" | "bold" = "medium";
    export let variant: "filled-surface" | "glass-surface" = "filled-surface";
    export let disabled: boolean = false;
</script>

<button class="btn btn-sm font-{weight} variant-{variant}" disabled={disabled} on:click>
    {text}
</button>