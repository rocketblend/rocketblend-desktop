import type {  ToastSettings } from '@skeletonlabs/skeleton';

import { t } from '$lib/translations/translations';

export const changeDetectedToast: ToastSettings = {
    message: t.get('home.toast.changeDetected'),
    timeout: 2000
};