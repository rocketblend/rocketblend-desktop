import type { PageLoad } from './$types';

import { application } from '$lib/wailsjs/go/models';
import { ListProjects } from '$lib/wailsjs/go/application/Driver'

export const load : PageLoad = (async ({ url }) => {
    const query = url.searchParams.get('query') || '';
    const opts = application.ListProjectsOpts.createFrom({
        query: query,
    });

    return {
        query: query,
        projects: (await ListProjects(opts)).projects,
    };
})