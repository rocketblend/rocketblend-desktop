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

export function convertToEnum(value: string, enumType: { [key: number]: string }): any {
    const num = parseInt(value, 10);
    if (!isNaN(num) && num in enumType) {
        return num as keyof typeof enumType;
    }
    return null;
}

export function formatTime(dateInput: Date | string): string {
    let date: Date;

    if (typeof dateInput === 'string') {
        date = new Date(dateInput);
    } else if (dateInput instanceof Date) {
        date = dateInput;
    } else {
        return "Invalid input";
    }

    // Check if 'date' is a valid Date object
    if (isNaN(date.getTime())) {
        return "Invalid date";
    }

    const options: Intl.DateTimeFormatOptions = {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    };

    return date.toLocaleTimeString(undefined, options);
}