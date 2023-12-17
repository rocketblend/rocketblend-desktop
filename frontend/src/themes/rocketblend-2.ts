import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const rocketblend2: CustomThemeConfig = {
    name: 'rocketblend-2',
    properties: {
        // =~= Theme Properties =~=
        "--theme-font-family-base": `system-ui`,
        "--theme-font-family-heading": `system-ui`,
        "--theme-font-color-base": "0 0 0",
        "--theme-font-color-dark": "255 255 255",
        "--theme-rounded-base": "4px",
        "--theme-rounded-container": "8px",
        "--theme-border-base": "1px",
        // =~= Theme On-X Colors =~=
        "--on-primary": "0 0 0",
        "--on-secondary": "0 0 0",
        "--on-tertiary": "255 255 255",
        "--on-success": "0 0 0",
        "--on-warning": "0 0 0",
        "--on-error": "255 255 255",
        "--on-surface": "255 255 255",
        // =~= Theme Colors  =~=
        // primary | #198ae1 
        "--color-primary-50": "221 237 251", // #ddedfb
        "--color-primary-100": "209 232 249", // #d1e8f9
        "--color-primary-200": "198 226 248", // #c6e2f8
        "--color-primary-300": "163 208 243", // #a3d0f3
        "--color-primary-400": "94 173 234", // #5eadea
        "--color-primary-500": "25 138 225", // #198ae1
        "--color-primary-600": "23 124 203", // #177ccb
        "--color-primary-700": "19 104 169", // #1368a9
        "--color-primary-800": "15 83 135", // #0f5387
        "--color-primary-900": "12 68 110", // #0c446e
        // secondary | #edb035 
        "--color-secondary-50": "252 243 225", // #fcf3e1
        "--color-secondary-100": "251 239 215", // #fbefd7
        "--color-secondary-200": "251 235 205", // #fbebcd
        "--color-secondary-300": "248 223 174", // #f8dfae
        "--color-secondary-400": "242 200 114", // #f2c872
        "--color-secondary-500": "237 176 53", // #edb035
        "--color-secondary-600": "213 158 48", // #d59e30
        "--color-secondary-700": "178 132 40", // #b28428
        "--color-secondary-800": "142 106 32", // #8e6a20
        "--color-secondary-900": "116 86 26", // #74561a
        // tertiary | #a61ddd 
        "--color-tertiary-50": "242 221 250", // #f2ddfa
        "--color-tertiary-100": "237 210 248", // #edd2f8
        "--color-tertiary-200": "233 199 247", // #e9c7f7
        "--color-tertiary-300": "219 165 241", // #dba5f1
        "--color-tertiary-400": "193 97 231", // #c161e7
        "--color-tertiary-500": "166 29 221", // #a61ddd
        "--color-tertiary-600": "149 26 199", // #951ac7
        "--color-tertiary-700": "125 22 166", // #7d16a6
        "--color-tertiary-800": "100 17 133", // #641185
        "--color-tertiary-900": "81 14 108", // #510e6c
        // success | #1ddd70 
        "--color-success-50": "221 250 234", // #ddfaea
        "--color-success-100": "210 248 226", // #d2f8e2
        "--color-success-200": "199 247 219", // #c7f7db
        "--color-success-300": "165 241 198", // #a5f1c6
        "--color-success-400": "97 231 155", // #61e79b
        "--color-success-500": "29 221 112", // #1ddd70
        "--color-success-600": "26 199 101", // #1ac765
        "--color-success-700": "22 166 84", // #16a654
        "--color-success-800": "17 133 67", // #118543
        "--color-success-900": "14 108 55", // #0e6c37
        // warning | #dd601d 
        "--color-warning-50": "250 231 221", // #fae7dd
        "--color-warning-100": "248 223 210", // #f8dfd2
        "--color-warning-200": "247 215 199", // #f7d7c7
        "--color-warning-300": "241 191 165", // #f1bfa5
        "--color-warning-400": "231 144 97", // #e79061
        "--color-warning-500": "221 96 29", // #dd601d
        "--color-warning-600": "199 86 26", // #c7561a
        "--color-warning-700": "166 72 22", // #a64816
        "--color-warning-800": "133 58 17", // #853a11
        "--color-warning-900": "108 47 14", // #6c2f0e
        // error | #dd1d53 
        "--color-error-50": "250 221 229", // #fadde5
        "--color-error-100": "248 210 221", // #f8d2dd
        "--color-error-200": "247 199 212", // #f7c7d4
        "--color-error-300": "241 165 186", // #f1a5ba
        "--color-error-400": "231 97 135", // #e76187
        "--color-error-500": "221 29 83", // #dd1d53
        "--color-error-600": "199 26 75", // #c71a4b
        "--color-error-700": "166 22 62", // #a6163e
        "--color-error-800": "133 17 50", // #851132
        "--color-error-900": "108 14 41", // #6c0e29
        // surface | #1c1c1c 
        "--color-surface-50": "221 221 221", // #dddddd
        "--color-surface-100": "210 210 210", // #d2d2d2
        "--color-surface-200": "198 198 198", // #c6c6c6
        "--color-surface-300": "164 164 164", // #a4a4a4
        "--color-surface-400": "96 96 96", // #606060
        "--color-surface-500": "28 28 28", // #1c1c1c
        "--color-surface-600": "25 25 25", // #191919
        "--color-surface-700": "21 21 21", // #151515
        "--color-surface-800": "17 17 17", // #111111
        "--color-surface-900": "14 14 14", // #0e0e0e
	},
    properties_dark: {
        "--color-surface-900": "00 00 00", // #000000
    }
}