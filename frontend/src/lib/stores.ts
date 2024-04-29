import { writable, get as getStoreValue } from 'svelte/store';
import type { ProjectIdStore, LogEvent, LogStore, OperationStore, PackageStore } from "$lib/types";
import type { types } from '$lib/wailsjs/go/models';

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

export const createOperationStore = (): OperationStore => {
    const { subscribe, set, update } = writable<types.Operation[]>([]);

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

export const createPackageStore = (): PackageStore => {
    const { subscribe, set, update } = writable<types.Package[]>([]);

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

export const getLogStore = () => logStore;

export const getSelectedProjectStore = () => selectedProjectStore;