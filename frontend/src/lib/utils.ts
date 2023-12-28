import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime';
import { v4 as uuidv4 } from 'uuid';

const heartbeatTimeout = 2000; // 2 seconds

export function cancellableOperationWithHeartbeat<T>(operation: (opID: string, ...args: any[]) => Promise<T>, ...args: any[]): [Promise<T | null>, () => void] {
    const opID = uuidv4();
    let cancelled = false;
    let heartbeatTimer: NodeJS.Timeout;

    const resetHeartbeatTimer = () => {
        clearTimeout(heartbeatTimer);
        heartbeatTimer = setTimeout(() => {
            if (!cancelled) {
                console.error("Operation timed out: " + opID);
                // Handle timeout, e.g., cancel operation, show error message, etc.
            }
        }, heartbeatTimeout);
    };

    // Start listening for heartbeats
    const cancelHeartbeatListener = EventsOn("operationHeartBeat", (data: string) => {
        if (data === opID) {
            resetHeartbeatTimer();
        }
    });

    resetHeartbeatTimer(); // Start the timer

    const operationPromise = new Promise<T | null>((resolve, reject) => {
        operation(opID, ...args).then(result => {
            clearTimeout(heartbeatTimer);
            cancelHeartbeatListener();
            if (!cancelled) {
                resolve(result);
            }
        }).catch(error => {
            clearTimeout(heartbeatTimer);
            cancelHeartbeatListener();
            if (!cancelled) {
                reject(error);
            }
        });
    });

    const cancel = () => {
        cancelled = true;
        clearTimeout(heartbeatTimer);
        cancelHeartbeatListener();
        EventsEmit("cancelOperation", opID);
    };

    return [operationPromise, cancel];
}