<script lang="ts">
    export let title: string;
    export let heading: "h1" | "h2" | "h3" | "h4" | "h5" | "h6" = "h6";
    export let weight: "medium" | "bold" = "bold";
</script>

<div class="{heading} font-{weight}">{title}</div>