import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime';
import { v4 as uuidv4 } from 'uuid';

const heartbeatTimeout = 2000; // 2 seconds

export function cancellableOperationWithHeartbeat<T>(operation: (opID: string, ...args: any[]) => Promise<T>, ...args: any[]): [Promise<T | null>, () => void] {
    const opID = uuidv4();
    let cancelled = false;
    let heartbeatTimer: NodeJS.Timeout;
    let rejectOperation: (reason?: any) => void;

    // console.log("Starting operation: " + opID);

    const resetHeartbeatTimer = () => {
        clearTimeout(heartbeatTimer);
        heartbeatTimer = setTimeout(() => {
            if (!cancelled) {
                console.error("Operation timed out: " + opID);
                cancel();
                rejectOperation(new Error("Operation timed out: " + opID));
            }
        }, heartbeatTimeout);
    };

    const cancelHeartbeatListener = EventsOn("operationHeartBeat", (data: string) => {
        // console.log("Received heartbeat: " + data);
        if (data === opID) {
            resetHeartbeatTimer();
        }
    });

    resetHeartbeatTimer();

    const operationPromise = new Promise<T | null>((resolve, reject) => {
        rejectOperation = reject;

        operation(opID, ...args).then(result => {
            // console.log("Operation completed: " + opID);
            clearTimeout(heartbeatTimer);
            cancelHeartbeatListener();
            if (!cancelled) {
                resolve(result);
            }
        }).catch(error => {
            if (!cancelled) {
                // console.log("Operation failed: " + opID);
                clearTimeout(heartbeatTimer);
                cancelHeartbeatListener();
                reject(error);
            }
        });
    });

    const cancel = () => {
        //console.log("Cancelled operation: " + opID);
        cancelled = true;
        clearTimeout(heartbeatTimer);
        cancelHeartbeatListener();
        EventsEmit("cancelOperation", opID);
    };

    return [operationPromise, cancel];
}