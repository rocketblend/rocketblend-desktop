<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { tableA11y } from './actions.js';

    import type { CssClasses, Tab } from "@skeletonlabs/skeleton";
    import type { TableSource, TableRow, TableColumn } from './types.js';

    import IconArrowUpFill from '~icons/ri/arrow-up-s-fill';
    import IconArrowDownFill from '~icons/ri/arrow-down-s-fill';

    const dispatch = createEventDispatcher();

    export let source: TableSource;
    export let selected: string[] = [];
    export let multiple: boolean = false;
    export let interactive = false;

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

    function toggleSort(column: TableColumn): void {
        // Clone the head array to trigger reactivity
        let newHead: TableColumn[] = source.head.map(h => {
            // Reset sort direction for other columns
            if (h.label !== column.label) {
                return { ...h, sortDirection: null };
            } 
            // Update the sort direction for the clicked column
            else {
                return {
                    ...h,
                    sortDirection: h.sortDirection === 'asc' ? 'desc' : 'asc'
                };
            }
        });

        // Update the source with the new head array
        source = { ...source, head: newHead };

        // Emit an event with the updated sorting parameters
        dispatch('sortChanged', { key: column.label, direction: column.sortDirection });
    }

    function onRowClick(clickedRow: TableRow): void {
        if (!interactive) return;

        const selectedIndex = selected.indexOf(clickedRow.id);
        if (multiple) {
            if (selectedIndex >= 0) {
                selected = selected.filter(id => id !== clickedRow.id);
            } else {
                selected = [...selected, clickedRow.id];
            }
        } else {
            selected = selectedIndex >= 0 ? [] : [clickedRow.id];
        }
    }

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
                    <th
                        class="{regionHeadCell} {heading.sortable ? 'sortable cursor-pointer' : ''}"
                        role="columnheader"
                        on:click={heading.sortable ? () => toggleSort(heading) : null}
                        tabindex={heading.sortable ? 0 : -1}
                        aria-label={heading.sortable ? `Sort by ${heading.label}` : heading.label}
                    >
                        <div class="inline-flex justify-center items-center space-x-2">
                            <div>{heading.label}</div>
                            {#if heading.sortable}
                                <div class="sort-indicator">
                                    {#if heading.sortDirection === 'asc'}
                                        <IconArrowUpFill />
                                    {:else if heading.sortDirection === 'desc'}
                                        <IconArrowDownFill />
                                    {/if}
                                </div>
                            {/if}
                        </div>  
                    </th>
                {/each}
            </tr>
        </thead>
        <tbody class="table-body {regionBody}">
            {#each source.body as row, rowIndex}
                <tr
                    class="{classesRow}"
                    class:table-row-checked={selected.includes(row.id)}
                    on:click={() => onRowClick(row)}
                    aria-rowindex={rowIndex + 1}
                    tabindex={interactive ? 0 : -1}>
                    {#each row.data as cell, cellIndex}
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