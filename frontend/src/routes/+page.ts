import type { PageLoad } from './$types';
import { FindAllProjects } from '$lib/wailsjs/go/application/Driver'

export const load : PageLoad = (async ({ url }) => {
    const query = url.searchParams.get('query') || '';

    return {
        query: query,
        projects: await FindAllProjects(query)
    };
})