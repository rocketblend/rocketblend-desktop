<script lang="ts">
    import FormatBytes from '$lib/components/ui/format/format-bytes.svelte';

    export let currentBytes: number;
    export let totalBytes: number;
</script>

<div class="flex items-center gap-2 whitespace-nowrap">
    <FormatBytes value={currentBytes}/>
    <div> / </div>
    <FormatBytes value={totalBytes}/>
</div>