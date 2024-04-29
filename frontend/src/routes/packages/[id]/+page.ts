import type { PageLoad } from './$types';
import { error} from '@sveltejs/kit';

import { GetPackage } from '$lib/wailsjs/go/application/Driver'
import { application } from '$lib/wailsjs/go/models';

export const load: PageLoad = async ({ params }) => {
    const opts = application.GetPackageOpts.createFrom({
        id: params.id,
    });

    const result = await GetPackage(opts);

    if (!result || result.package === undefined) {
        throw error(404, {
            message: 'Package not found'
        });
    }

    return {
        label: result.package.name+" "+result.package.tag,
        package: result.package
    };
};