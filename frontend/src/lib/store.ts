import { writable, get as getStoreValue } from 'svelte/store';
import type { LogEvent } from './types';

function createProjectIdStore() {
    const { subscribe, set, update } = writable<string[]>([]);

    return {
        subscribe,

        // Remove a specific ID from the array
        remove: (id: string) => update(ids => {
            return ids.filter(existingId => existingId !== id);
        }),

        // Get the latest element of the array
        latest: () => {
            const ids = getStoreValue({ subscribe });
            return ids.length > 0 ? ids[ids.length - 1] : null;
        },

        // Get the current state of the array
        get: () => getStoreValue({ subscribe }),

        // Replace the entire array (useful for initialization or clearing)
        set: (ids: string[]) => set(ids),

        // Clear the entire array
        clear: () => set([])
    };
}

export const selectedProjectIds = createProjectIdStore();