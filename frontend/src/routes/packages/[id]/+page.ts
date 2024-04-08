import type { PageLoad } from './$types';
import { error} from '@sveltejs/kit';
import { GetPackage } from '$lib/wailsjs/go/application/Driver'

export const load: PageLoad = async ({ params }) => {
    const result = await GetPackage(params.id);

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