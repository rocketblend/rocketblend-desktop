<script lang="ts">
    import { writable } from "svelte/store";
    import { tableA11y } from './actions.js';

    import type { CssClasses } from "@skeletonlabs/skeleton";
    import type { TableSource } from './types.js';

    export let source: TableSource;
    export let interactive = false;

    // Style Props...
    export let element: CssClasses = 'table';
    export let text: CssClasses = '';
    export let color: CssClasses = '';
    export let regionHead: CssClasses = '';
    export let regionHeadCell: CssClasses = '';
    export let regionBody: CssClasses = '';
    export let regionCell: CssClasses = '';
    export let regionFoot: CssClasses = '';
    export let regionFootCell: CssClasses = '';
    export let regionRow: CssClasses = '';

    // Bindable Property for Selected Rows
    export let selectedRows = writable<number[]>([]);

    // Row Click Handler
    function onRowClick(rowIndex: number): void {
        if (!interactive) return;
        selectedRows.update(current => {
            const index = current.indexOf(rowIndex);
            if (index >= 0) {
                return current.filter(i => i !== rowIndex); // Remove if already selected
            } else {
                return [...current, rowIndex]; // Add to selection if not
            }
        });
    }

    // Reactive variables for classes
    $: classesBase = `${$$props.class || ''}`;
    $: classesTable = `${element} ${text} ${color}`;
    $: classesRow = `${regionRow}`;
</script>

<div class="table-container {classesBase}">
    <table
        class="{classesTable}"
        class:table-interactive={interactive}
        role={interactive ? "grid" : "table"}
        use:tableA11y>
        <thead class="table-head {regionHead}">
            <tr>
                {#each source.head as heading}
                    <th class="{regionHeadCell}" role="columnheader">{@html heading}</th>
                {/each}
            </tr>
        </thead>
        <tbody class="table-body {regionBody}">
            {#each source.body as row, rowIndex}
                <tr
                    class="{classesRow}"
                    class:table-row-checked={$selectedRows.includes(rowIndex)}
                    on:click={() => onRowClick(rowIndex)}
                    aria-rowindex={rowIndex + 1}
                    tabindex={interactive ? 0 : -1}>
                    {#each row as cell, cellIndex}
                        <td
                            class="{regionCell}"
                            role="gridcell"
                            aria-colindex={cellIndex + 1}
                            tabindex={cellIndex === 0 && interactive ? 0 : -1}>
                            {@html Number(cell) === 0 ? cell : (cell ? cell : '-')}
                        </td>
                    {/each}
                </tr>
            {/each}
        </tbody>
        {#if source.foot}
            <tfoot class="table-foot {regionFoot}">
                <tr>
                    {#each source.foot as cell}
                        <td class="{regionFootCell}">{@html cell}</td>
                    {/each}
                </tr>
            </tfoot>
        {/if}
    </table>
</div>