<script lang="ts">
    import { LongRunningOperation } from '$lib/wailsjs/go/application/Driver'
    import { cancellableOperationWithHeartbeat } from '$lib/utils';

    import { t } from '$lib/translations/translations';

    let operationStatus = 'Not started';
    let operationPromise: Promise<void | null>;
    let cancelOperation: () => void;

    function startOperation() {
        operationStatus = 'Running...';
        const [opPromise, cancelFunc] = cancellableOperationWithHeartbeat<void>(LongRunningOperation, 15000);
        operationPromise = opPromise;
        cancelOperation = cancelFunc;

        operationPromise.then(() => {
            if (operationStatus !== 'Cancelled') {
                operationStatus = 'Completed';
            }
        }).catch(error => {
            if (operationStatus !== 'Cancelled') {
                operationStatus = `Failed: ${error}`;
            }
        });
    }

    function cancel() {
        if (cancelOperation) {
            cancelOperation();
            operationStatus = 'Cancelled';
        }
    }
</script>

<div class="flex flex-col card p-2 space-y-2">
    <button class="btn variant-filled" on:click={startOperation} disabled={operationStatus === 'Running...'}>
        {$t('home.debug.operation.start')}
    </button>
    
    <button class="btn variant-filled" on:click={cancel} disabled={operationStatus !== 'Running...'}>
        {$t('home.debug.operation.cancel')}
    </button>
    
    <p>{$t('home.debug.operation.status')}: {operationStatus}</p>
</div>