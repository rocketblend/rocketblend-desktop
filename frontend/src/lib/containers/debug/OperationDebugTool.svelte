<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { LongRunningOperation, CancelOperation, ListOperations, GetOperation } from '$lib/wailsjs/go/application/Driver';
    import type { operationservice } from '$lib/wailsjs/go/models';
    import { EventsOn, EventsOff } from '$lib/wailsjs/runtime';
    
    import { t } from '$lib/translations/translations';
    import { writable } from 'svelte/store';

    let cooldown = false;
    let operations = writable<Array<operationservice.Operation>>([]);

    function startOperation() {
        cooldown = true;
        setTimeout(() => {
            cooldown = false;
        }, 5000); // 5 seconds cooldown

        LongRunningOperation().then(response => {
            fetchOperations();
            GetOperation(response).then(operationDetails => {
                console.log('Operation started:', operationDetails);
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
            operations.set([...result]);
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

        EventsOn('storeEvent', (data: { id: string, type: number, indexType: string }) => {
            if (data.indexType === "operation") {
                GetOperation(data.id).then(operation => {
                    console.log('Operation updated', operation);
                    fetchOperations();
                });
            }
        });
    });

    onDestroy(() => {
        EventsOff('storeEvent');
    });
</script>

<div class="flex flex-col card p-2 space-y-2">
    <button class="btn variant-filled" on:click={startOperation} disabled={cooldown}>
        {$t('home.debug.operation.start')}
    </button>
    <hr>
    <ul class="space-y-1">
        {#each $operations as operation (operation.id)}
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