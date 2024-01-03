import type { ToastSettings, ToastStore } from '@skeletonlabs/skeleton';

import { t } from '$lib/translations/translations';

import type { LogStore, LogEvent } from '$lib/types';
import {debounce} from '$lib/utils';
import { resourcePath } from '$lib/components/utils';
import { EventsOn, EventsOff, EventsEmit } from '$lib/wailsjs/runtime';

export const setupGlobalEventListeners = (logStore: LogStore, toastStore: ToastStore) => {
    const changeDetectedDebounce = debounce(() => {
        const changeDetectedToast: ToastSettings = {
            message: t.get('home.toast.changeDetected'),
            timeout: 3000
        };
        toastStore.trigger(changeDetectedToast);
    }, 1000);


    // Setup debug log listener
    EventsOn('debug.log', (data: LogEvent) => {
        logStore.add(data);
    });

    // Setup application argument listener
    EventsOn('application.argument', (data: { args: string[] }) => {
        if (data.args && data.args.length !== 0) {
            const applicationArgumentToast: ToastSettings = {
                message: `Args: ${data.args.join(', ')}`,
                timeout: 5000,
            };

            toastStore.trigger(applicationArgumentToast);
        }
    });

    // Setup search store listener
    EventsOn('searchstore.insert', (data: { id: string, indexType: string }) => {
        changeDetectedDebounce();
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
};

export const tearDownGlobalEventListeners = () => {
    // Remove log stream listener
    EventsOff('debug.log');

    // Remove search store listener
    EventsOff('searchstore.insert');

    // Remove launch arguments listener
    EventsOff('application.argument');
}