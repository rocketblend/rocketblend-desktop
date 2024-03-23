import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';
import { superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';

import { desktopGeneralForm } from './(components)/forms';

export const load: PageLoad = async ({ params }) => {
    // const id = parseInt(params.id);
  
    // const request = await fetch(
    //   `https://jsonplaceholder.typicode.com/users/${id}`
    // );
    // if (request.status >= 400) throw error(request.status);
  
    // const userData = await request.json();
    const form = await superValidate(zod(desktopGeneralForm));
  
    return { form };
};