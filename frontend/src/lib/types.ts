export type RadioOption = {
    value: number;
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