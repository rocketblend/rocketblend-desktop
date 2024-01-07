<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { Chart } from 'frappe-charts';
    import type { AxisOptions, BarOptions, LineOptions, TooltipOptions, ChartData, ChartType } from 'frappe-charts';

    export let data: ChartData;

    export let title: string = '';
    export let type: ChartType = 'line';
    export let height: number = 300;
    export let animate: boolean = true;
    export let axisOptions: AxisOptions = {};
    export let barOptions: BarOptions = {};
    export let lineOptions: LineOptions = {};
    export let tooltipOptions: TooltipOptions=  {};
    export let colors: string[] = [];
    export let valuesOverPoints: number = 0;
    export let isNavigable: boolean = false;
    export let maxSlices: number = 3;
    export let radius: number = 0;
    export let discreteDomains: number = 0;

    let chart: Chart | null = null;
    let chartRef: HTMLElement;

    function ifChartThen<T extends (...args: any[]) => any>(fn: T): (...funcArgs: Parameters<T>) => ReturnType<T> | undefined {
        return (...args) => {
            if (chart) {
                return fn(...args);
            }
        }
    }

    export const addDataPoint = ifChartThen((label: string, valueFromEachDataset: number[], index: number) => chart?.addDataPoint(label, valueFromEachDataset, index));
    export const removeDataPoint = ifChartThen((index: number) => chart?.removeDataPoint(index));
    export const exportChart = ifChartThen(() => chart?.export());
    const updateChart = ifChartThen((newData: typeof data) => chart?.update(newData));
    $: updateChart(data);

    onMount(() => {
        console.log(data)
        chart = new Chart(chartRef, {
            data: data,
            title,
            type,
            height,
            animate,
            colors,
            axisOptions,
            barOptions,
            lineOptions,
            tooltipOptions,
            valuesOverPoints,
            isNavigable,
            maxSlices,
            discreteDomains,
            radius,
        });
    });

    onDestroy(() => {
        chart = null;
    });
</script>

<div class="overflow-x-auto" bind:this={chartRef} on:data-select>
</div>