import type { ToastSettings, ToastStore } from '@skeletonlabs/skeleton';
import { EventsOn, EventsOff, EventsEmit } from '$lib/wailsjs/runtime';
import { t } from '$lib/translations/translations';
import type { LogStore, LogEvent } from '$lib/types';

export function setupGlobalEventListeners(logStore: LogStore, toastStore: ToastStore) {
    // Setup log stream listener
    EventsOn('logStream', (data: LogEvent) => {
        logStore.add(data);
    });

    // Setup launch arguments listener
    EventsOn('launchArgs', (data: { args: string[] }) => {
        if (data.args && data.args.length !== 0) {
            const launchToast: ToastSettings = {
                message: `Args: ${data.args.join(', ')}`,
                timeout: 5000,
            };

            toastStore.trigger(launchToast);
        }
    });

    EventsOn('storeEvent', (data: { id: string, type: number, indexType: string }) => {
        console.log('storeEvent', data);
        if (data) {
            const storeToast: ToastSettings = {
                message: `Store event: ${data.id} ${data.type} ${data.indexType}`,
                timeout: 5000,
            };

            toastStore.trigger(storeToast);
        }
    });

    // Emit a ready event for the backend to listen for.
    EventsEmit('ready');
    
    // Trigger an initial toast
    const initialToast: ToastSettings = {
        message: t.get('home.greeting'),
        background: 'bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500 text-white',
        timeout: 5000,
    };
    toastStore.trigger(initialToast);
}

export function tearDownGlobalEventListeners() {
    // Remove log stream listener
    EventsOff('logStream');

    // Remove launch arguments listener
    EventsOff('launchArgs');
}