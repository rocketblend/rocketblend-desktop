<script lang="ts">
    import IconBox2Fill from '~icons/ri/box-2-fill';
    import IconLoopRightFill from '~icons/ri/loop-right-fill';
    import IconAddBoxFill from '~icons/ri/add-box-fill';

    export let title: string;
    
    // Event handlers as props
    export let onRefresh: () => void;
    export let onAdd: () => void;
</script>
  
<div class="inline-flex items-center align-center space-x-2 text-surface-200">
    <IconBox2Fill />
    <h5 class="font-bold flex-grow">{title}</h5>
    <button type="button" class="btn p-0 text-surface-200" on:click={onRefresh}>
        <IconLoopRightFill class="text-md mt-1"/>
    </button>
    <button type="button" class="btn p-0 text-surface-200" on:click={onAdd}>
        <IconAddBoxFill class="text-xl mt-1"/>
    </button>
</div>