export interface TableRow {
    id: string; // Unique identifier for each row
    data: string[]; // Data for the row
}

export interface TableSource {
    /** The formatted table heading values. */
    head: string[];
    /** The formatted table body values. */
    body: TableRow[];
    /** The formatted table footer values. */
    foot?: string[];
}