export type Loading = 'eager' | 'lazy';

export type ImageDetails = {
    src: string;
    alt: string;
    class: string;
    loading: Loading;
};