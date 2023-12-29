import { writable, get as getStoreValue } from 'svelte/store';
import type { ProjectIdStore, LogEvent, LogStore } from "$lib/types";

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
        addLog: (logItem: LogEvent) => update(logs => [...logs, logItem]),
        clearLogs: () => set([]),
    };
};

const logStore = createLogStore();

export const getLogStore = () => logStore;

export const getSelectedProjectStore = () => selectedProjectStore;