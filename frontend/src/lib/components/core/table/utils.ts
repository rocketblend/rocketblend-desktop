export function getNestedObjectValue(obj: any, key: string): any {
    let value = obj;
    key.split('.').forEach((k) => {
        if(value != null) {
            value = value[k];
        }
    });
    return value;
}

export function tableSourceMapper(source: any[], keys: string[]): any[] {
    return source.map((row) => {
        const mappedRow: any = {};
        keys.forEach((key) => (mappedRow[key] = getNestedObjectValue(row, key)));
        return mappedRow;
    });
}

/** Map an array of objects to an array of values. */
export function tableSourceValues(source: any[]): any[] {
    return source.map((row) => Object.values(row));
}

/** Sets object order and returns values. */
export function tableMapperValues(source: any[], keys: string[]): any[] {
    return tableSourceValues(tableSourceMapper(source, keys));
}