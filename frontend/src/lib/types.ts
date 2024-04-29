import type { types } from '$lib/wailsjs/go/models';

export type RadioOption = {
    value: number | string;
    display: string;
};

export type MediaInfo = {
    id: string;
    title: string;
    src: string;
    alt?: string;
};

export enum SortBy {
    Name,
    File,
    Build
}

export enum DisplayType {
    Table,
    Gallery
}

export type Option = {
    value: number;
    display: string;
};

export type OptionGroup = {
    label: string;
    display: string;
    options: Option[];
};

export type LogEvent = {
    level: string,
    message: string,
    time: Date,
    fields: { [key: string]: string }
}

export type ProjectIdStore = {
    subscribe: (run: (value: string[]) => void) => () => void;
    remove: (id: string) => void;
    latest: () => string | null;
    get: () => string[];
    set: (ids: string[]) => void;
    clear: () => void;
};

export type LogStore = {
    subscribe: (run: (value: LogEvent[]) => void) => () => void;
    add: (logItem: LogEvent) => void;
    clear: () => void;
};

export type OperationStore = {
    subscribe: (this: void, run: import("svelte/store").Subscriber<types.Operation[]>, invalidate?: (value?: types.Operation[]) => void) => import("svelte/store").Unsubscriber;
    set: (operations: types.Operation[]) => void;
    add: (operation: types.Operation) => void;
    clear: () => void;
}

export type PackageStore = {
    subscribe: (this: void, run: import("svelte/store").Subscriber<types.Package[]>, invalidate?: (value?: types.Package[]) => void) => import("svelte/store").Unsubscriber;
    set: (packages: types.Package[]) => void;
    add: (pack: types.Package) => void;
    clear: () => void;
}