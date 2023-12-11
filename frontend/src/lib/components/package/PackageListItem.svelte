<script lang="ts">
    export let name: string;
    export let reference: string;
</script>

<div class="flex gap-2">
    <!-- Render package details -->
    <div class="flex-grow w-2">
        <div class="text-sm font-medium">{name}</div>
        <div class="text-sm">{reference}</div>
    </div>
</div>