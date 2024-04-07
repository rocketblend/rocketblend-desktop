<script lang="ts">
    import IconAddFill from '~icons/ri/add-fill';
    import IconSubtractFill from '~icons/ri/subtract-fill';

    export let assigned: boolean = false;
    export let variantFrom: string = 'secondary';
    export let variantTo: string = 'tertiary';
    export let active: boolean = false;
    export let rounded: boolean = true;

    function handleAction() {
        assigned = !assigned;
    }

    function handleUserInteraction(event: KeyboardEvent | MouseEvent) {
        if (event.type === 'click' || (event instanceof KeyboardEvent && event.key === 'Enter')) {
            handleAction();
        }
    }
</script>

<button
    class="btn h-full bg-gradient-to-br variant-gradient-{variantFrom}-{variantTo} {rounded ? 'rounded' : ''} p-1 text-token"
    on:click|stopPropagation={handleUserInteraction}
    on:keydown|stopPropagation={handleUserInteraction}
>
    {#if active}
        {#if assigned}
            <IconSubtractFill class="w-5 h-5" />
        {:else}
            <IconAddFill class="w-5 h-5" />
        {/if}
    {/if}
</button>