import type { PageLoad } from "./$types";

export const load : PageLoad = (async ({ url }) => {
    return {
        label: "Getting Started",
    }
})