import type { PageLoad } from './$types';

import { error} from '@sveltejs/kit';

import { GetProject } from '$lib/wailsjs/go/application/Driver'

export const load: PageLoad = async ({ params }) => {
    const result = await GetProject(params.id);

    if (!result || result.project === undefined) {
        throw error(404, {
            message: 'Project not found'
        });
    }

    return {
        label: result.project.name,
        project: result.project
    };
};