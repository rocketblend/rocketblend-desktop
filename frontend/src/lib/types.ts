export type RadioOption = {
    value: number;
    key: string;
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

export type OperationEntry = {
    key: string;
    cancel: () => void;
};

export type CancellableOperationsStore = {
    subscribe: (this: void, run: import("svelte/store").Subscriber<OperationEntry[]>, invalidate?: (value?: OperationEntry[]) => void) => import("svelte/store").Unsubscriber;
    add: (entry: OperationEntry) => void;
    cancel: (key: string) => void;
    cancelAll: () => void;
};