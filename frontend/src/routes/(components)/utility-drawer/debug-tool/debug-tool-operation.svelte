<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    import { ProgressRadial, getToastStore } from '@skeletonlabs/skeleton';

    import { t } from '$lib/translations/translations';
    import { LongRunningOperation, CancelOperation, ListOperations, GetOperation } from '$lib/wailsjs/go/application/Driver';
    import type { operationservice } from '$lib/wailsjs/go/models';
    import { EventsOn } from '$lib/wailsjs/runtime';
    import { createOperationStore } from '$lib/stores';
    import { debounce } from '$lib/utils';
	import { EVENT_DEBOUNCE } from '$lib/events';

    const operationStore = createOperationStore();
    const toastStore = getToastStore();
    const debounceFetchOperations = debounce(fetchOperations, EVENT_DEBOUNCE);

    let cooldown = false;
    let cancelListener: () => void;

    function startOperation() {
        cooldown = true;
        setTimeout(() => {
            cooldown = false;
        }, 2000); // 2 seconds cooldown

        LongRunningOperation().then(response => {
            fetchOperations();
            GetOperation(response).then(operationDetails => {
                console.log('Operation started', operationDetails);
            }).catch(error => {
                console.log('Error fetching operation details:', error);
            });

        }).catch(error => {
            console.log(`Operation start failed: ${error}`);
        });
    }

    function cancelOperation(id: string) {
        CancelOperation(id).then(() => {
            console.log(`Operation cancelled: ${id}`);
        }).catch(error => {
            console.log(`Operation cancel failed: ${error}`);
        });
    }

    function fetchOperations() {
        ListOperations().then(result => {
            operationStore.set([...result]);
        }).catch(error => {
            console.log(`Error fetching operations: ${error}`);
        });
    }

    function getStatusText(operation: operationservice.Operation): string {
        if (operation.completed && operation.error) {
            return $t('home.debug.operation.cancelled');
        } else if (operation.error) {
            return $t('home.debug.operation.failed');
        } else if (operation.completed) {
            return $t('home.debug.operation.completed');
        } else {
            return $t('home.debug.operation.running');
        }
    }

    onMount(() => {
        fetchOperations();

        cancelListener = EventsOn('searchstore.insert', (data: { id: string, indexType: string }) => {
            if (data.indexType === "operation") {
                debounceFetchOperations();
            }
        });
    });

    onDestroy(() => {
        if (cancelListener) {
            cancelListener();
        }
    });
</script>

<div class="flex flex-col card p-2 space-y-2">
    <button class="btn variant-filled" on:click={startOperation} disabled={cooldown}>
        {$t('home.debug.operation.start')}
    </button>
    <hr>
    <ul class="space-y-1">
        {#each $operationStore as operation}
            <li class="flex flex-col card p-2 space-y-2 text-sm">
                <div>
                    <div>ID: {operation.id}</div>
                    <div>Status: {getStatusText(operation)}</div>
                    {#if operation.error }
                        <div>Error: {operation.error}</div>
                    {/if}
                </div>
                {#if !operation.completed}
                    <button class="btn variant-ghost-warning w-full" on:click={() => cancelOperation(operation.id.toString())}>
                        <div class="flex justify-center items-center space-x-2">
                            <span>Cancel</span>
                            <ProgressRadial width="w-4" stroke={40} strokeLinecap="round"/>
                        </div>
                    </button>
                {/if}
            </li>
        {/each}
    </ul>
</div>