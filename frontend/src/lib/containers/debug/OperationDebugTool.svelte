<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    import type {  ToastSettings } from '@skeletonlabs/skeleton';
    import { getToastStore } from '@skeletonlabs/skeleton';

    import { LongRunningOperation, CancelOperation, ListOperations, GetOperation } from '$lib/wailsjs/go/application/Driver';
    import type { operationservice } from '$lib/wailsjs/go/models';
    import { EventsOn, EventsOff } from '$lib/wailsjs/runtime';

    import { t } from '$lib/translations/translations';
	import { getOperationStore } from '$lib/stores';

    const operationStore = getOperationStore();
    const toastStore = getToastStore();

    let cooldown = false;
    let cancelListener: () => void;

    function startOperation() {
        cooldown = true;
        setTimeout(() => {
            cooldown = false;
        }, 5000); // 5 seconds cooldown

        LongRunningOperation().then(response => {
            fetchOperations();
            GetOperation(response).then(operationDetails => {
                // const toast: ToastSettings = {
                //     message: t.get('home.operation.started.message', { id: operationDetails.id.toString() }),
                // };

                // toastStore.trigger(toast);
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
                fetchOperations();
                // GetOperation(data.id).then(operation => {
                //     console.log('Operation updated', operation);
                //     fetchOperations();
                // });
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
            <li class="flex flex-col card p-2 space-y-2">
                <div>
                    ID: {operation.id}, Status: {getStatusText(operation)}, Error: {operation.error || 'None'}
                </div>
                {#if !operation.completed}
                    <button class="btn variant-ghost-warning w-full" on:click={() => cancelOperation(operation.id.toString())}>Cancel</button>
                {/if}
            </li>
        {/each}
    </ul>
</div>