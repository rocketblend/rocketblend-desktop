import { invalidateAll } from '$app/navigation';

import type { ToastSettings, ToastStore } from '@skeletonlabs/skeleton';

import { t } from '$lib/translations/translations';
import type { LogStore, LogEvent } from '$lib/types';
import {debounce} from '$lib/utils';
import { EventsOn, EventsOff, EventsEmit } from '$lib/wailsjs/runtime';

export const EVENT_DEBOUNCE = 250;
export const SEARCH_STORE_INSERT_CHANNEL = 'store.insert';
export const DEBUG_LOG_CHANNEL = 'debug.log';
export const APPLICATION_ARGUMENT_CHANNEL = 'application.argument';

export const setupGlobalEventListeners = (logStore: LogStore, toastStore: ToastStore) => {
    const changeDetectedDebounce = debounce(() => {
        invalidateAll(); // Invalidate all routes (force a re-render of the app)        
        const changeDetectedToast: ToastSettings = {
            message: t.get('home.toast.changeDetected'),
            timeout: 3000
        };
        toastStore.trigger(changeDetectedToast);
    }, EVENT_DEBOUNCE - 50);

    // Setup debug log listener
    EventsOn(DEBUG_LOG_CHANNEL, (data: LogEvent) => {
        logStore.add(data);
    });

    // Setup application argument listener
    EventsOn(APPLICATION_ARGUMENT_CHANNEL, (data: { args: string[] }) => {
        if (data.args && data.args.length !== 0) {
            const applicationArgumentToast: ToastSettings = {
                message: `Args: ${data.args.join(', ')}`,
                timeout: 5000,
            };

            toastStore.trigger(applicationArgumentToast);
        }
    });

    // Setup search store listener
    EventsOn(SEARCH_STORE_INSERT_CHANNEL, (data: { id: string, indexType: string }) => {
        if (data.indexType === 'operation') {
            return;
        }

        changeDetectedDebounce();
    });

    // Emit a ready event for the backend to listen for.
    EventsEmit('ready');
    
    // Trigger an initial toast
    const initialToast: ToastSettings = {
        message: t.get('home.greeting'),
        background: 'bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500 text-white',
        timeout: 10000,
    };
    toastStore.trigger(initialToast);
};

export const tearDownGlobalEventListeners = () => {
    // Remove log stream listener
    EventsOff(DEBUG_LOG_CHANNEL);

    // Remove search store listener
    EventsOff(SEARCH_STORE_INSERT_CHANNEL);

    // Remove launch arguments listener
    EventsOff(APPLICATION_ARGUMENT_CHANNEL);
}