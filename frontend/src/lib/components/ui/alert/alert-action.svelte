<script lang="ts">
    export let text: string;
    export let weight: "medium" | "bold" = "medium";
    export let variant: "filled-surface" | "glass-surface" = "filled-surface";
</script>

<button class="btn btn-sm font-{weight} variant-{variant}" on:click>
    {text}
</button>