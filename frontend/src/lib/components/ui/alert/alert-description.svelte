<script lang="ts">
    export let message: string;
    export let size: "sm" | "md" | "lg" = "sm";
</script>

<p class="text-{size}">{message}</p>