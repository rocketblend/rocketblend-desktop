import type { PageLoad } from './$types';
import { FindProjectByKey } from '$lib/wailsjs/go/application/Driver'

export const load = (async ({ params }) => {
    return {
        projects: await FindProjectByKey(params.key)
    };
}) satisfies PageLoad;