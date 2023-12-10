<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import SearchInput from '$lib/components/core/input/SearchInput.svelte';
    import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';

    export let form: HTMLFormElement;
    export let searchQuery: string;
    export let displayType: string;

    const dispatch = createEventDispatcher();

    function handleChange() {
        dispatch('change')
    }
</script>

<form bind:this={form} class="inline-flex space-x-4 w-full" data-sveltekit-keepfocus>
    <div class="flex-grow">
        <SearchInput name="query" value={searchQuery} on:input={handleChange} placeholder="Search" debounceDelay={500} />
    </div>
    <RadioGroup>
        <RadioItem name="display" group={displayType} on:change={handleChange} value="table">Table</RadioItem>
        <RadioItem name="display" group={displayType} on:change={handleChange} value="gallery">Gallery</RadioItem>
    </RadioGroup>
    <button type="submit" class="hidden">Search</button>
</form>