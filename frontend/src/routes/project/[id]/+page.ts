import type { PageLoad } from './$types';
import { FindProjectByID } from '$lib/wailsjs/go/application/Driver'

export const load = (async ({ params }) => {
    return {
        project: await FindProjectByID(params.id)
    };
}) satisfies PageLoad;