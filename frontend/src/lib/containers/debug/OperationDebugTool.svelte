<script lang="ts">
    import { onDestroy } from 'svelte';

    import { LongRunningWithOperation } from '$lib/wailsjs/go/application/Driver';
    import { cancellableOperationWithHeartbeat } from '$lib/utils';
    import { t } from '$lib/translations/translations';

    enum OperationStatus {
        NotStarted,
        Running,
        Completed,
        Cancelled,
        Failed,
    }

    let operationStatus = OperationStatus.NotStarted;
    let operationPromise: Promise<void | null>;
    let cancelOperation: () => void;

    function startOperation() {
        operationStatus = OperationStatus.Running;
        const [opPromise, cancelFunc] = cancellableOperationWithHeartbeat<void>(LongRunningWithOperation, 15000);
        operationPromise = opPromise;
        cancelOperation = cancelFunc;

        operationPromise.then(() => {
            if (operationStatus !== OperationStatus.Cancelled) {
                operationStatus = OperationStatus.Completed;
            }
        }).catch(error => {
            if (operationStatus !== OperationStatus.Cancelled) {
                operationStatus = OperationStatus.Failed;
                console.error(`Operation failed: ${error}`); // Handle the error appropriately
            }
        });
    }

    function cancel() {
        if (cancelOperation) {
            cancelOperation();
            operationStatus = OperationStatus.Cancelled;
        }
    }

    function getStatusText(status: OperationStatus): string {
        switch (status) {
            case OperationStatus.NotStarted:
                return $t('home.debug.operation.notStarted');
            case OperationStatus.Running:
                return $t('home.debug.operation.running');
            case OperationStatus.Completed:
                return $t('home.debug.operation.completed');
            case OperationStatus.Cancelled:
                return $t('home.debug.operation.cancelled');
            case OperationStatus.Failed:
                return $t('home.debug.operation.failed');
            default:
                return '';
        }
    }

    onDestroy(() => {
        if (operationStatus === OperationStatus.Running) {
            cancel();
        }
    });
</script>

<div class="flex flex-col card p-2 space-y-2">
    <button class="btn variant-filled" on:click={startOperation} disabled={operationStatus === OperationStatus.Running}>
        {$t('home.debug.operation.start')}
    </button>
    
    <button class="btn variant-filled" on:click={cancel} disabled={operationStatus !== OperationStatus.Running}>
        {$t('home.debug.operation.cancel')}
    </button>
    
    <p>{$t('home.debug.operation.status')}: {getStatusText(operationStatus)}</p>
</div>