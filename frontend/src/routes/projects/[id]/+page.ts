import type { PageLoad } from './$types';

import { error} from '@sveltejs/kit';

import { GetProject } from '$lib/wailsjs/go/application/Driver'
import { application } from '$lib/wailsjs/go/models';

export const load: PageLoad = async ({ params }) => {
    const opts = application.GetPackageOpts.createFrom({
        id: params.id,
    });

    const result = await GetProject(opts);

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