
declare namespace svelteHTML {
    import { ChartDataPoint } from 'frappe-charts';

    interface HTMLAttributes<T> {
        'on:data-select'?: (event: CustomEvent<ChartDataPoint>) => void;
    }
}