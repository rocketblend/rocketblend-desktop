export function debounce(fn: Function, delay: number) {
    let timer: ReturnType<typeof setTimeout>;

    return function (...args: any[]) {
        clearTimeout(timer);
        timer = setTimeout(() => fn(...args), delay);
    };
}

export function resourcePath(path: string | undefined) {
    if (path && path != "") {
        return `/system/${path}`;
    }

    return "";
};

export const videoExtensions = ['.webm', '.mp4', '.ogg'];
export function isVideo(path: string) {
    return videoExtensions.some(ext => path.endsWith(ext));
}