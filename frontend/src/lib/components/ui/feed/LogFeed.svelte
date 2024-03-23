<script lang="ts">
    import { tick } from 'svelte';
	import type { LogEvent } from "$lib/types";
	import { formatTime } from "../../utils";

    export let feed: LogEvent[] = [];

    let elemFeed: HTMLElement;
    let isFirstMount = true;

    function scrollChatBottom(): void {
        if (!elemFeed) {
            return;
        }

        const behavior = isFirstMount ? 'instant' : 'smooth';
        elemFeed.scrollTo({ top: elemFeed.scrollHeight, behavior });

        if (isFirstMount) {
            isFirstMount = false;
        }
    }

    function logLevelClass(level: string): string {
        switch (level.toLowerCase()) {
            case 'error':
                return 'text-error-500-400-token';
            case 'warning':
                return 'text-warning-500-400-token';
            case 'info':
                return 'text-success-500-400-token';
            default:
                return '';
        }
    }

    function fieldValueClass(key: string): string {
        return key.toLowerCase() === 'err' ? 'text-error-500-400-token' : '';
    }

    $: if (feed.length) {
        tick().then(() => {
            scrollChatBottom();
        });
    }
</script>

<section bind:this={elemFeed} class="h-full overflow-y-auto space-y-1 text-sm lowercase text-surface-900-50-token">
	{#each feed as log}
        <div class="w-full">
            <span class="text-surface-600-300-token">{formatTime(log.time)}</span>
            <span> | </span>
            <span class="uppercase {logLevelClass(log.level)}">{log.level}</span>
            <span> | </span>
            <span class="font-bold">{log.message}</span>
            <span> | </span>
            {#each Object.entries(log.fields) as [key, value]}
                <span class="text-primary-500-400-token">{key}=</span><span class="{fieldValueClass(key)}">{value + " "}</span>
            {/each}
        </div>
	{/each}
</section>
					