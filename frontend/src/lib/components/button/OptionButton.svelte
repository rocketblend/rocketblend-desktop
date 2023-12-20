<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { ListBox, ListBoxItem, popup } from '@skeletonlabs/skeleton';
    import type { PopupSettings } from '@skeletonlabs/skeleton';
    import IconListUnordered from '~icons/ri/list-unordered';
    import type { OptionGroup } from '$lib/types';

    export let optionsGroups: OptionGroup[];
    export let selectedOptions: Record<string, number>;

    const dispatch = createEventDispatcher();

    const popupSettings: PopupSettings = {
        event: 'focus-click',
        target: 'OptionPopup',
        placement: 'bottom',
        closeQuery: '.listbox-item'
    };

    function handleOptionChange(groupLabel: string, value: number) {
        selectedOptions[groupLabel] = value;
        dispatch('optionChange', { groupLabel, value });
    }
</script>

<button class="btn justify-between variant-ghost-surface" use:popup={popupSettings}>
        <slot name="buttonContent">
            <!-- Default content if no slot content is provided -->
            Options
        </slot>
</button>
<div class="card w-48 p-2 variant-filled-surface rounded z-50" data-popup={popupSettings.target}>
    {#each optionsGroups as group (group.label)}
        <ListBox padding="px-2 py-2">
            <div class="pb-2 px-2 pt-1">
                <span class="text-xs font-bold text-surface-700-200-token">{group.display}</span>
            </div>
            {#each group.options as option (option.value)}
                <ListBoxItem bind:group={selectedOptions[group.label]} name={group.label} value={option.value} on:change={() => handleOptionChange(group.label, option.value)}>
                    {option.display}
                </ListBoxItem>
            {/each}
        </ListBox>
    {/each}
    <div class="arrow variant-filled-surface" />
</div>