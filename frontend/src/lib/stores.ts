import { writable, get as getStoreValue } from 'svelte/store';
import type { ProjectIdStore, LogEvent, LogStore, OperationStore, PackageStore } from "$lib/types";
import type { operationservice, pack } from '$lib/wailsjs/go/models';

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

const createOperationStore = (): OperationStore => {
    const { subscribe, set, update } = writable<operationservice.Operation[]>([]);

    return {
        subscribe,
        set: (operations) => {
            set(operations);
        },
        add: (operation) => {
            update(operations => [...operations, operation]);
        },
        clear: () => {
            set([]);
        },
    };
};

const operationStore = createOperationStore();

const createPacakgeStore = (): PackageStore => {
    const { subscribe, set, update } = writable<pack.Package[]>([]);

    return {
        subscribe,
        set: (pack) => {
            set(pack);
        },
        add: (pack) => {
            update(packages => [...packages, pack]);
        },
        clear: () => {
            set([]);
        },
    };
};

const packageStore = createPacakgeStore();

export const getPackageStore = () => packageStore;

export const getOperationStore = () => operationStore;

export const getLogStore = () => logStore;

export const getSelectedProjectStore = () => selectedProjectStore;