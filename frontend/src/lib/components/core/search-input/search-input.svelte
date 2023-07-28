<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { debounce } from '$lib/components/utils';
  
    export let placeholder: string = '';
    export let debounceDelay: number = 250;
  
    let searchText: string = '';
  
    const dispatch = createEventDispatcher();
  
    const handleInput = debounce((value: string) => {
        dispatch('search', value);
    }, debounceDelay);
  </script>
  
  <div class="relative">
    <label for="search" class="sr-only">{placeholder}</label>
    <input
      bind:value={searchText}
      on:input={({ target: { value } }) => handleInput(value)}
      id="search"
      type="text"
      class="input"
      {placeholder}
    />
</div>