import { writable, get as getStoreValue } from 'svelte/store';
import type { ProjectIdStore, LogEvent, LogStore, OperationEntry, CancellableOperationsStore } from "$lib/types";

const MAX_LOGS = 1000;

const createSelectedProjectStore = (): ProjectIdStore => {
    const { subscribe, set, update } = writable<string[]>([]);

    return {
        subscribe,
        remove: (id: string) => update(ids => ids.filter(existingId => existingId !== id)),
        latest: () => {
            const ids = getStoreValue({ subscribe });
            return ids.length > 0 ? ids[ids.length - 1] : null;
        },
        get: () => getStoreValue({ subscribe }),
        set: (ids: string[]) => set(ids),
        clear: () => set([])
    };
};

const selectedProjectStore = createSelectedProjectStore();

const createLogStore = (): LogStore => {
    const { subscribe, set, update } = writable<LogEvent[]>([]);

    return {
        subscribe,
        add: (logItem: LogEvent) => {
            update(logs => {
                const updatedLogs = [...logs, logItem];

                if (updatedLogs.length > MAX_LOGS) {
                    return updatedLogs.slice(-MAX_LOGS);
                }

                return updatedLogs;
            });
        },
        clear: () => set([]),
    };
};

const logStore = createLogStore();

const createCancellableOperationsStore = (): CancellableOperationsStore => {
    const { subscribe, update } = writable<OperationEntry[]>([]);

    return {
        subscribe,
        add: (entry: OperationEntry) => {
            update(operations => [...operations, entry]);
        },
        cancel: (key: string) => {
            update(operations => {
                const operationIndex = operations.findIndex(op => op.key === key);
                if (operationIndex !== -1) {
                    operations[operationIndex].cancel();
                    operations.splice(operationIndex, 1);
                }
                return operations;
            });
        },
        cancelAll: () => {
            const operations = getStoreValue({ subscribe });
            operations.forEach(entry => entry.cancel());
            update(() => []);
        }
    };
};

const cancellableOperationsStore = createCancellableOperationsStore();

export const getCancellableOperationsStore = () => cancellableOperationsStore;

export const getLogStore = () => logStore;

export const getSelectedProjectStore = () => selectedProjectStore;