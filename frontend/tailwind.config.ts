import { join } from 'path';
import type { Config } from 'tailwindcss';

import { skeleton } from '@skeletonlabs/tw-plugin';

import { rocketblend } from './src/themes/rocketblend';

const config = {
	darkMode: 'class',
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		join(require.resolve(
			'@skeletonlabs/skeleton'),
			'../**/*.{html,js,svelte,ts}'
		)
	],
    safelist: [
        {
          pattern: /h-\d+/,
        },
        {
          pattern: /w-\d+/,
        },
        {
          pattern: /variant-gradient-.*/,
        },
    ],
	theme: {
		extend: {},
	},
	plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
		skeleton({
            themes: {
                custom: [ rocketblend ],
                // preset: [
                //     "skeleton",
                //     "modern",
                //     "crimson",
                //     "wintry",
                //     "hamlindigo"
                // ]
            }
        })
	]
} satisfies Config;

export default config;