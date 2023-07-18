import type { PageLoad } from './$types';
import { FindAllProjects, Greet } from '$lib/wailsjs/go/application/Driver'

export const load = (async ({ params }) => {
    return {
        greeting: await Greet("Test"),
        projects: await FindAllProjects()
    };

}) satisfies PageLoad;