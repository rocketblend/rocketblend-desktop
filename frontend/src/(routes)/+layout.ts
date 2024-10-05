import { locale, loadTranslations } from '$lib/translations/translations';
import type { LayoutLoad } from './$types';

import { getSelectedProjectStore } from '$lib/stores';
import { GetDetails, GetPreferences, GetProject, ListPackages } from '$lib/wailsjs/go/application/Driver'
import { application, types } from '$lib/wailsjs/go/models';

export const ssr = false;
export const prerender = "auto";

export const load: LayoutLoad = async ({ url, depends }) => {
    const { pathname } = url;
    const defaultLocale = 'en'; // get from cookie / user session etc...
    const initLocale: string = locale.get() || defaultLocale;

    const selectedProjectStore = getSelectedProjectStore();
    const selectedId = selectedProjectStore.latest();

    await loadTranslations(initLocale, pathname); // keep this just before the `return`

    depends('app:layout');

    let project: types.Project | undefined;
    if (selectedId) {
        await GetProject(application.GetPackageOpts.createFrom({
            id: selectedId,
        })).then((result) => {
            project = result.project;
        }).catch(() => {
            selectedProjectStore.clear();
        });
    }

    return {
        showBreadcrumb: true,
        selectedProject: project,
        details: GetDetails(),
        preferences: await GetPreferences(),
    }
};