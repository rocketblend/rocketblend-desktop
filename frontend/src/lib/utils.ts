import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime';
import { v4 as uuidv4 } from 'uuid';

const heartbeatTimeout = 10000; // 10 seconds

export function cancellableOperationWithHeartbeat<T>(operation: (opid: string, ...args: any[]) => Promise<T>, ...args: any[]): [Promise<T | null>, () => void] {
    const opid = uuidv4();
    let cancelled = false;
    let heartbeatTimer: NodeJS.Timeout;
    let rejectOperation: (reason?: any) => void;

    // console.log("Starting operation: " + opid);

    const resetHeartbeatTimer = () => {
        clearTimeout(heartbeatTimer);
        heartbeatTimer = setTimeout(() => {
            if (!cancelled) {
                // console.error("Operation timed out: " + opid);
                cancel();
                rejectOperation(new Error("Operation timed out: " + opid));
            }
        }, heartbeatTimeout);
    };

    const cancelHeartbeatListener = EventsOn("operationHeartBeat", (data: string) => {
        // console.log("Received heartbeat: " + data);
        if (data === opid) {
            resetHeartbeatTimer();
        }
    });

    resetHeartbeatTimer();

    const operationPromise = new Promise<T | null>((resolve, reject) => {
        rejectOperation = reject;

        operation(opid, ...args).then(result => {
            // console.log("Operation completed: " + opid);
            clearTimeout(heartbeatTimer);
            cancelHeartbeatListener();
            if (!cancelled) {
                resolve(result);
            }
        }).catch(error => {
            if (!cancelled) {
                // console.log("Operation failed: " + opid);
                clearTimeout(heartbeatTimer);
                cancelHeartbeatListener();
                reject(error);
            }
        });
    });

    const cancel = () => {
        //console.log("Cancelled operation: " + opid);
        cancelled = true;
        clearTimeout(heartbeatTimer);
        cancelHeartbeatListener();
        EventsEmit("cancelOperation", opid);
    };

    return [operationPromise, cancel];
}