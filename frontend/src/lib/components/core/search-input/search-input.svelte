<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { debounce } from '$lib/components/utils';
  
  let clazz: string = '';
  export { clazz as class };

  export let placeholder: string = '';
  export let value: string = '';
  export let name: string = 'search';
  export let debounceDelay: number = 250;

  const dispatch = createEventDispatcher();

  const processInput = debounce((event: Event) => {
    const input = event.target as HTMLInputElement;
      dispatch('input', input.value);
  }, debounceDelay);
</script>

<div class="relative">
  <label for="search" class="sr-only">{placeholder}</label>
  <input
    bind:value={value}
    on:input={processInput}
    id="search"
    type="text"
    class="input {clazz}"
    autocomplete="off"
    {placeholder}
    {name}
  />
</div>