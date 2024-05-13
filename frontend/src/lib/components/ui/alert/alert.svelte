<script lang="ts">
    import IconInfoFill from '~icons/ri/file-info-fill';

    import AlertTitle from './alert-title.svelte';

    export let variant: "ghost-surface" | "ghost-primary" | "ghost-warning" = "ghost-surface";
</script>

<aside class="alert variant-{variant}">
    <slot name="icon">
        <IconInfoFill />
    </slot>
    <div class="alert-message  w-full">
        <slot name="title">
            <AlertTitle title="Info" />
        </slot>
        <slot></slot>
    </div>
    <div class="alert-actions">
        <slot name="actions">
        </slot>
    </div>
</aside>