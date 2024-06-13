<script lang="ts">
    export let size: "sm" | "md" | "lg" = "sm";
</script>

<p class="text-{size}">
    <slot></slot>
</p>