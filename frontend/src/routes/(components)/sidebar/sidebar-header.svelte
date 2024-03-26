<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    import IconBox2Fill from '~icons/ri/box-2-fill';
    import IconLoopRightFill from '~icons/ri/loop-right-fill';
    import IconAddBoxFill from '~icons/ri/add-box-fill';

    export let title: string;

    const dispatch = createEventDispatcher();
    
    function refresh() {
        dispatch('refresh');
    }

    function add() {
        dispatch('add');
    }
</script>
  
<div class="inline-flex items-center space-x-2 text-surface-200">
    <IconBox2Fill />
    <h5 class="font-bold flex-grow">{title}</h5>
    <button type="button" aria-label="Refresh" class="btn p-0 text-surface-200" on:click={refresh}>
        <IconLoopRightFill class="text-md mt-1"/>
    </button>
    <button type="button" aria-label="Add" class="btn p-0 text-surface-200" on:click={add}>
        <IconAddBoxFill class="text-xl mt-1"/>
    </button>
</div>