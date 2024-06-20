import type { PageLoad } from './$types';

import { error} from '@sveltejs/kit';

import { GetProject, ListPackages } from '$lib/wailsjs/go/application/Driver'
import { application, types } from '$lib/wailsjs/go/models';

export const load: PageLoad = async ({ params }) => {
    const result = await GetProject(application.GetPackageOpts.createFrom({
        id: params.id,
    }))

    if (!result || result.project === undefined) {
        throw error(404, {
            message: 'Project not found'
        });
    }

    const packages = await ListPackages(application.ListPackagesOpts.createFrom({
        references: result.project.dependencies?.map(dep => dep.reference) || [],
    }));

    return {
        label: result.project.name,
        project: result.project,
        dependencies: packages.packages || [],
    };
};