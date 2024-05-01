export function debounce(fn: Function, delay: number) {
    let timer: ReturnType<typeof setTimeout>;

    return function (...args: any[]) {
        clearTimeout(timer);
        timer = setTimeout(() => fn(...args), delay);
    };
}

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

const parseAndValidateDate = (dateInput: Date | string): Date | null => {
    let date: Date;

    if (typeof dateInput === 'string') {
        date = new Date(dateInput);
    } else if (dateInput instanceof Date) {
        date = dateInput;
    } else {
        return null;
    }

    if (isNaN(date.getTime())) {
        return null;
    }

    return date;
}

export const formatDateTime = (dateInput: Date | string): string => {
    const date = parseAndValidateDate(dateInput);
    if (!date) return "Invalid input or date";

    const options: Intl.DateTimeFormatOptions = {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    };

    return date.toLocaleString(undefined, options);
};

export const formatTime = (dateInput: Date | string): string => {
    const date = parseAndValidateDate(dateInput);
    if (!date) return "Invalid input or date";

    const options: Intl.DateTimeFormatOptions = {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    };

    return date.toLocaleTimeString(undefined, options);
};