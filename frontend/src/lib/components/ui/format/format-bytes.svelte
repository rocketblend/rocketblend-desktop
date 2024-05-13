<script lang="ts">
    export let value: number;
    export let suffix: string = '';
  
    function formatBytes(bytes: number): string {
        const units = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    
        if (bytes === 0) {
            return `0 Bytes${suffix}`;
        }
    
        const index = Math.floor(Math.log(bytes) / Math.log(1024));
        const convertedValue = bytes / Math.pow(1024, index);
        return `${convertedValue.toFixed(2)} ${units[index]}${suffix}`;
    }
</script>
  
<div class="text-sm">{formatBytes(value)}</div>