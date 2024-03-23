<script lang="ts">
    import { ProgressBar } from '@skeletonlabs/skeleton';

    export let value: number = 0;
    export let min: number = 0;
    export let max: number = 100;
    export let rounded: boolean = true;
    export let toFixed: number = 0;

    $: percentage = ((value - min) / (max - min)) * 100;
    $: formattedPercentage = Math.min(100, Math.max(0, percentage)); // Clamp between 0 and 100
</script>

<div class="flex items-center gap-2">
    <ProgressBar value={value} min={min} max={max} rounded={rounded ? "rounded" : ""} />
    <slot>
        <div class="text-surface-800-100-token text-sm">{formattedPercentage.toFixed(toFixed)}%</div>
    </slot>
</div>