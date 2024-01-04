<script lang="ts">
    import { tick, createEventDispatcher } from 'svelte';

    import { debounce } from '$lib/utils';

    const dispatch = createEventDispatcher();
    interface Option { label: string; value: string; }
    interface InputTypeState {
        text: boolean;
        number: boolean;
        textarea: boolean;
        select: boolean;
    }

    const getInputTypeState = (type: string): InputTypeState => ({
        text: type === 'text',
        number: type === 'number',
        textarea: type === 'textarea',
        select: type === 'select'
    });

    export let value: string = '';
    export let type: 'text' | 'number' | 'textarea' | 'select' = 'text';
    export let placeholder: string = '';
    export let labelClasses: string = '';
    export let inputClasses: string = 'input';
    export let rows: number = 2;
    export let cols: number = 20;
    export let options: Option[] = [];

    let editing: boolean = false;
    let inputEl: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement;
    let label: string;
    let selectedIndex: number = options.findIndex(o => o.value === value);
    let currentInputType: InputTypeState = getInputTypeState(type);

    let lastConfirmedValue: string = value;
    let isKeyPressHandled = false;

    const computeLabel = (): string => {
        if (currentInputType.text || currentInputType.number || currentInputType.textarea) {
            return value || placeholder;
        } else if (currentInputType.select) {
            return selectedIndex === -1 ? placeholder : options[selectedIndex].label;
        }
        
        return '';
    };

    $: currentInputType = getInputTypeState(type);
    $: label = computeLabel();

    const focusInput = async () => {
        await tick();
        inputEl?.focus();
    };

    const toggleEditing = () => {
        if (!editing) {
            lastConfirmedValue = value;
        }

        editing = !editing;

        if (editing) {
            focusInput();
        } else if (value !== lastConfirmedValue) {
            label = computeLabel();
            dispatch('change', value);
        }
    };

    const handleInput = (event: Event) => {
        value = (event.target as HTMLInputElement).value;
    };

    const handleKeyPress = (event: KeyboardEvent) => {
        console.log(event.key);
        if (event.key === 'Enter' && editing) {
            toggleEditing();
            isKeyPressHandled = true;
        }
    };

    const handleKeyDown = (event: KeyboardEvent) => {
        if (event.key === 'Escape') {
            value = lastConfirmedValue;
            toggleEditing();
            isKeyPressHandled = true;
        }
    };

    const handleBlur = () => {
        if (!isKeyPressHandled) {
            toggleEditing();
            dispatch('blur', value);
        }

        isKeyPressHandled = false;
    };

    const handleChange = (event: Event) => {
        const target = event.target as HTMLSelectElement;
        selectedIndex = options.findIndex(option => option.value === target.value);
        value = options[selectedIndex]?.value || '';
    };
</script>

<div>
    {#if editing}
        <!-- Split inputs to handle types correctly with two-way binding -->
        {#if currentInputType.text}
            <input
                class={inputClasses}
                bind:this={inputEl}
                type="text"
                bind:value
                placeholder={placeholder}
                on:input={handleInput}
                on:blur={handleBlur}
                on:keypress={handleKeyPress}
                on:keydown={handleKeyDown}
            />
        {:else if currentInputType.number}
            <input
                class={inputClasses}
                bind:this={inputEl}
                type="number"
                bind:value={value}
                placeholder={placeholder}
                on:input={handleInput}
                on:blur={handleBlur}
                on:keypress={handleKeyPress}
            />
        {:else if currentInputType.textarea}
            <textarea
                class={inputClasses}
                bind:this={inputEl}
                bind:value
                placeholder={placeholder}
                rows={rows}
                cols={cols}
                on:input={handleInput}
                on:blur={handleBlur}
            />
        {:else if currentInputType.select}
            <select
                class={inputClasses}
                bind:this={inputEl}
                bind:value
                on:change={handleChange}
                on:blur={handleBlur}
            >
                {#if placeholder}
                    <option disabled>{placeholder}</option>
                {/if}
                {#each options as option}
                    <option value={option.value}>
                        {option.label}
                    </option>
                {/each}
            </select>
        {/if}
    {:else}
        <div
            class="inline-flex justify-end items-center space-x-2 {labelClasses}"
            on:click={toggleEditing}
            on:keypress={handleKeyPress}
            tabindex="0"
            role="button"
            aria-label="Interactive label"
            aria-expanded={editing}
        >
            <div>
                {label}
            </div>
            <slot></slot>
            <slot name="selectCaret">{#if currentInputType.select}<span>&#9660;</span>{/if}</slot>
        </div>
    {/if}
</div>