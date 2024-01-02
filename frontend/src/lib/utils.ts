import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime';
import { v4 as uuidv4 } from 'uuid';

const heartbeatTimeout = 10000; // 10 seconds

export function cancellableOperationWithHeartbeat<T>(operation: (cid: string, ...args: any[]) => Promise<T>, ...args: any[]): [Promise<T | null>, () => void] {
    const cid = uuidv4();
    let cancelled = false;
    let heartbeatTimer: NodeJS.Timeout;
    let rejectOperation: (reason?: any) => void;

    // console.log("Starting operation: " + cid);

    const resetHeartbeatTimer = () => {
        clearTimeout(heartbeatTimer);
        heartbeatTimer = setTimeout(() => {
            if (!cancelled) {
                // console.error("Operation timed out: " + cid);
                cancel();
                rejectOperation(new Error("Operation timed out: " + cid));
            }
        }, heartbeatTimeout);
    };

    const cancelHeartbeatListener = EventsOn("requestHeartBeat", (data: string) => {
        // console.log("Received heartbeat: " + data);
        if (data === cid) {
            resetHeartbeatTimer();
        }
    });

    resetHeartbeatTimer();

    const operationPromise = new Promise<T | null>((resolve, reject) => {
        rejectOperation = reject;

        operation(cid, ...args).then(result => {
            // console.log("Operation completed: " + cid);
            clearTimeout(heartbeatTimer);
            cancelHeartbeatListener();
            if (!cancelled) {
                resolve(result);
            }
        }).catch(error => {
            if (!cancelled) {
                // console.log("Operation failed: " + cid);
                clearTimeout(heartbeatTimer);
                cancelHeartbeatListener();
                reject(error);
            }
        });
    });

    const cancel = () => {
        //console.log("Cancelled operation: " + cid);
        cancelled = true;
        clearTimeout(heartbeatTimer);
        cancelHeartbeatListener();
        EventsEmit("operation.cancel", cid);
    };

    return [operationPromise, cancel];
}

export function debounce<T extends (...args: any[]) => any>(func: T, wait: number): (...args: Parameters<T>) => void {
    let timeout: ReturnType<typeof setTimeout> | null = null;

    return function(...args: Parameters<T>): void {
        const later = () => {
            timeout = null;
            func(...args);
        };

        if (timeout !== null) {
            clearTimeout(timeout);
        }
        timeout = setTimeout(later, wait);
    };
}