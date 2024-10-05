<script lang="ts">
    import { onDestroy } from 'svelte';

    import { LongRunningRequestWithCancellation } from '$lib/wailsjs/go/application/Driver';
    
    import { cancellableOperationWithHeartbeat } from '$lib/utils';
    import { t } from '$lib/translations/translations';

    enum RequestStatus {
        NotStarted,
        Running,
        Completed,
        Cancelled,
        Failed,
    }

    let requestStatus = RequestStatus.NotStarted;
    let requestPromise: Promise<void | null>;
    let cancelOperation: () => void;

    function startRequest() {
        requestStatus = RequestStatus.Running;
        const [opPromise, cancelFunc] = cancellableOperationWithHeartbeat<void>(LongRunningRequestWithCancellation, 15000);
        requestPromise = opPromise;
        cancelOperation = cancelFunc;

        requestPromise.then(() => {
            if (requestStatus !== RequestStatus.Cancelled) {
                requestStatus = RequestStatus.Completed;
            }
        }).catch(error => {
            if (requestStatus !== RequestStatus.Cancelled) {
                requestStatus = RequestStatus.Failed;
                console.error(`Request failed: ${error}`); // Handle the error appropriately
            }
        });
    }

    function cancel() {
        if (cancelOperation) {
            cancelOperation();
            requestStatus = RequestStatus.Cancelled;
        }
    }

    function getStatusText(status: RequestStatus): string {
        switch (status) {
            case RequestStatus.NotStarted:
                return $t('home.debug.request.notStarted');
            case RequestStatus.Running:
                return $t('home.debug.request.running');
            case RequestStatus.Completed:
                return $t('home.debug.request.completed');
            case RequestStatus.Cancelled:
                return $t('home.debug.request.cancelled');
            case RequestStatus.Failed:
                return $t('home.debug.request.failed');
            default:
                return '';
        }
    }

    onDestroy(() => {
        if (requestStatus === RequestStatus.Running) {
            cancel();
        }
    });
</script>

<div class="flex flex-col card p-2 space-y-2">
    <button class="btn variant-filled" on:click={startRequest} disabled={requestStatus === RequestStatus.Running}>
        {$t('home.debug.request.start')}
    </button>
    
    <button class="btn variant-filled" on:click={cancel} disabled={requestStatus !== RequestStatus.Running}>
        {$t('home.debug.request.cancel')}
    </button>
    
    <p>{$t('home.debug.request.status')}: {getStatusText(requestStatus)}</p>
</div>