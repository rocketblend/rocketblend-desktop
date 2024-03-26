export function getNestedObjectValue(obj: any, key: string): any {
    let value = obj;
    key.split('.').forEach((k) => {
        if(value != null) {
            value = value[k];
        }
    });
    return value;
}

/** Map an array of objects to an array of values. */
export function tableSourceValues(source: any[]): any[] {
    return source.map((row) => Object.values(row));
}

export function tableSourceMapper(source: any[], idKey: string, keys: string[]): any[] {
    return source.map((row) => {
        const id = getNestedObjectValue(row, idKey);
        const data = keys.map(key => getNestedObjectValue(row, key));

        return { id, data };
    });
}

export function tableMapperValues(source: any[], idKey: string, keys: string[]): any[] {
    return tableSourceMapper(source, idKey, keys);
}