<script lang="ts">
    import { tick } from 'svelte';
	import type { LogEvent } from "$lib/types";
	import { formatTime } from "../utils";

    export let feed: LogEvent[] = [];

    let elemFeed: HTMLElement;

    function scrollChatBottom(behavior?: ScrollBehavior): void {
        if (!elemFeed) {
            return;
        }

        elemFeed.scrollTo({ top: elemFeed.scrollHeight, behavior });
    }

    $: if (feed) {
        tick().then(() => {
            scrollChatBottom('smooth');
        });
    }
</script>

<section bind:this={elemFeed} class="h-full overflow-y-auto space-y-1 text-sm lowercase text-surface-900-50-token">
	{#each feed as log}
        <div class="space-x-1 w-full">
            <span>{formatTime(log.time)}</span>
            <span>|</span>
            <span class="uppercase">{log.level}</span>
            <span>|</span>
            <span class="font-bold">{log.message}</span>
            <span>|</span>
            {#each Object.entries(log.fields) as [key, value]}
                <span>
                    <span class="text-primary-500-400-token">{key}=</span><span>{value}</span>
                </span>
            {/each}
        </div>
	{/each}
</section>
					