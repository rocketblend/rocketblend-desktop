<script lang="ts">
    export let label = '';
    export let variant = 'ghost';
    export let rounded = true;
</script>

<div class={`badge variant-${variant} ${rounded ? 'rounded' : ''}`}>
    {label}
</div>