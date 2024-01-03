<script lang="ts">
    import { tick, createEventDispatcher } from 'svelte';
    
    const dispatch = createEventDispatcher();

    interface Option {
        label: string;
        value: string;
    }

    export let value: string = '';
    export let type: 'text' | 'number' | 'textarea' | 'select' = 'text';
    export let placeholder: string = '';
    export let labelClasses: string = '';
    export let inputClasses: string = '';
    export let rows: number = 2;
    export let cols: number = 20;
    export let options: Option[] = [];

    let editing: boolean = false;
    let inputEl: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement;
    let label: string;
    let selectedIndex: number = options.findIndex(o => o.value === value);
    
    $: isText = type === 'text';
    $: isNumber = type === 'number';
    $: isTextArea = type === 'textarea';
    $: isSelect = type === 'select';
    $: if (isNumber) {
        label = value === '' ? placeholder : value;
    } else if (isText || isTextArea) {
        label = value ? value : placeholder;
    } else {
        label = selectedIndex === -1 ? placeholder : options[selectedIndex].label;
    }
    
    const toggle = async () => {
        editing = !editing;

        if (editing) {
            await tick();
            inputEl.focus();
        }
    };
    
    const handleInput = (e: Event) => {
        const target = e.target as HTMLInputElement;
        value = isNumber ? target.value : target.value;
    };
    
    const handleEnter = (e: KeyboardEvent) => {
        if (e.key === 'Enter') inputEl.blur();
    };
    
    const handleBlur = () => {
        toggle();
        dispatch('blur', value);
    };

    const handleKeyPress = (e: KeyboardEvent) => {
        if (e.key === 'Enter' || e.key === ' ') {
            toggle();
        }
    };
    
    const handleChange = (e: Event) => {
        const target = e.target as HTMLSelectElement;
        selectedIndex = placeholder ? target.selectedIndex - 1 : target.selectedIndex;
        value = options[selectedIndex].value;
    };
</script>
    
{#if editing && (isText || isNumber)}
    <input
        class={inputClasses}
        bind:this={inputEl}
        {type}
        {value}
        {placeholder}
        on:input={handleInput}
        on:keyup={handleEnter}
        on:blur={handleBlur}
    >
{:else if editing && isTextArea}
    <textarea
        class={inputClasses}
        bind:this={inputEl}
        {placeholder}
        {value}
        {rows}
        {cols}
        on:input={handleInput}
        on:blur={handleBlur}
    />
{:else if editing && isSelect}
    <select
        class={inputClasses}
        bind:this={inputEl}
        on:change={handleChange}
        {value}
        on:blur={handleBlur}
    >
    {#if placeholder}
        <option selected value disabled>{placeholder}</option>
    {/if}
    {#each options as { label, value }, i}
        <option value={value}>
            {label}
        </option>
    {/each}
    </select>
{:else}
    <span
        class={labelClasses}
        on:click={toggle}
        on:keypress={handleKeyPress}
        tabindex="0"
        role="button"
        aria-label="Your descriptive label here"
    >
        {label}<slot name="selectCaret">{#if isSelect}<span>&#9660;</span>{/if}</slot>
    </span>
{/if}