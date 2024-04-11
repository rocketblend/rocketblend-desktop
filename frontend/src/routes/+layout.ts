import { locale, loadTranslations } from '$lib/translations/translations';
import type { LayoutLoad } from './$types';

import { getSelectedProjectStore } from '$lib/stores';
import { GetDetails, GetPreferences, GetProject } from '$lib/wailsjs/go/application/Driver'

export const ssr = false;
export const prerender = "auto";

export const load: LayoutLoad = async ({ url, depends }) => {
    const { pathname } = url;
    const defaultLocale = 'en'; // get from cookie / user session etc...
    const initLocale: string = locale.get() || defaultLocale;
    const selectedId = getSelectedProjectStore().latest();

    await loadTranslations(initLocale, pathname); // keep this just before the `return`

    depends('app:layout');

    return {
        showBreadcrumb: true,
        details: GetDetails(),
        selectedProject: selectedId ? await GetProject(selectedId) : undefined,
        preferences: await GetPreferences(),
    }
};