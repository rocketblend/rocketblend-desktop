import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const rocketblend2: CustomThemeConfig = {
    name: 'rocketblend-2',
    properties: {
        // =~= Theme Properties =~=
		"--theme-font-family-base": `system-ui`,
		"--theme-font-family-heading": `system-ui`,
		"--theme-font-color-base": "0 0 0",
		"--theme-font-color-dark": "255 255 255",
		"--theme-rounded-base": "2px",
		"--theme-rounded-container": "8px",
		"--theme-border-base": "1px",
		// =~= Theme On-X Colors =~=
		"--on-primary": "0 0 0",
		"--on-secondary": "255 255 255",
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
		// secondary | #aa26df 
		"--color-secondary-50": "242 222 250", // #f2defa
		"--color-secondary-100": "238 212 249", // #eed4f9
		"--color-secondary-200": "234 201 247", // #eac9f7
		"--color-secondary-300": "221 168 242", // #dda8f2
		"--color-secondary-400": "196 103 233", // #c467e9
		"--color-secondary-500": "170 38 223", // #aa26df
		"--color-secondary-600": "153 34 201", // #9922c9
		"--color-secondary-700": "128 29 167", // #801da7
		"--color-secondary-800": "102 23 134", // #661786
		"--color-secondary-900": "83 19 109", // #53136d
		// tertiary | #b23455 
		"--color-tertiary-50": "243 225 230", // #f3e1e6
		"--color-tertiary-100": "240 214 221", // #f0d6dd
		"--color-tertiary-200": "236 204 213", // #ecccd5
		"--color-tertiary-300": "224 174 187", // #e0aebb
		"--color-tertiary-400": "201 113 136", // #c97188
		"--color-tertiary-500": "178 52 85", // #b23455
		"--color-tertiary-600": "160 47 77", // #a02f4d
		"--color-tertiary-700": "134 39 64", // #862740
		"--color-tertiary-800": "107 31 51", // #6b1f33
		"--color-tertiary-900": "87 25 42", // #57192a
		// success | #84cc16 
		"--color-success-50": "237 247 220", // #edf7dc
		"--color-success-100": "230 245 208", // #e6f5d0
		"--color-success-200": "224 242 197", // #e0f2c5
		"--color-success-300": "206 235 162", // #ceeba2
		"--color-success-400": "169 219 92", // #a9db5c
		"--color-success-500": "132 204 22", // #84cc16
		"--color-success-600": "119 184 20", // #77b814
		"--color-success-700": "99 153 17", // #639911
		"--color-success-800": "79 122 13", // #4f7a0d
		"--color-success-900": "65 100 11", // #41640b
		// warning | #edb035 
		"--color-warning-50": "252 243 225", // #fcf3e1
		"--color-warning-100": "251 239 215", // #fbefd7
		"--color-warning-200": "251 235 205", // #fbebcd
		"--color-warning-300": "248 223 174", // #f8dfae
		"--color-warning-400": "242 200 114", // #f2c872
		"--color-warning-500": "237 176 53", // #edb035
		"--color-warning-600": "213 158 48", // #d59e30
		"--color-warning-700": "178 132 40", // #b28428
		"--color-warning-800": "142 106 32", // #8e6a20
		"--color-warning-900": "116 86 26", // #74561a
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