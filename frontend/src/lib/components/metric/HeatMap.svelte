<script lang="ts">
	import Chart from "$lib/components/core/chart/Chart.svelte";
    import type { HeatmapData, ChartType } from 'frappe-charts';

    const dateToUnix = (dateString: string | Date): number => {
        return new Date(dateString).getTime();
    };

    let data: HeatmapData = {
        dataPoints: {
            [dateToUnix('2024-01-01')]: 4,
            [dateToUnix('2024-01-02')]: 8,
            [dateToUnix('2024-01-03')]: 15,
            [dateToUnix('2024-01-04')]: 16,
            [dateToUnix('2024-01-05')]: 23,
            [dateToUnix('2024-01-06')]: 5,
            [dateToUnix('2024-01-12')]: 15,
            // ... more data points
        },
        start: new Date('2024-01-01'), // Optional: start date for the heatmap
        end: new Date('2024-12-31'),   // Optional: end date for the heatmap
    };

    const chartType = 'heatmap';
</script>

<Chart
    data={data}
    type={chartType}
    height={175}
    radius={1}
    discreteDomains={1}
/>