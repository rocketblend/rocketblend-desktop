import type { PageLoad } from './$types';
import { ListProjects, ListPackages } from '$lib/wailsjs/go/application/Driver'

export const load : PageLoad = (async ({ url }) => {
    const query = url.searchParams.get('query') || '';

    return {
        query: query,
        projects: (await ListProjects(query)).projects,
        packages: (await ListPackages(query)).packages
    };
})