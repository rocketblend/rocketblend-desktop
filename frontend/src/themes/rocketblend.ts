import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const rocketblend: CustomThemeConfig = {
    name: 'rocketblend',
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
		"--on-secondary": "0 0 0",
		"--on-tertiary": "255 255 255",
		"--on-success": "0 0 0",
		"--on-warning": "0 0 0",
		"--on-error": "255 255 255",
		"--on-surface": "255 255 255",
		// =~= Theme Colors  =~=
		// primary | #3e9ae0 
		"--color-primary-50": "226 240 250", // #e2f0fa
		"--color-primary-100": "216 235 249", // #d8ebf9
		"--color-primary-200": "207 230 247", // #cfe6f7
		"--color-primary-300": "178 215 243", // #b2d7f3
		"--color-primary-400": "120 184 233", // #78b8e9
		"--color-primary-500": "62 154 224", // #3e9ae0
		"--color-primary-600": "56 139 202", // #388bca
		"--color-primary-700": "47 116 168", // #2f74a8
		"--color-primary-800": "37 92 134", // #255c86
		"--color-primary-900": "30 75 110", // #1e4b6e
		// secondary | #e6b145 
		"--color-secondary-50": "251 243 227", // #fbf3e3
		"--color-secondary-100": "250 239 218", // #faefda
		"--color-secondary-200": "249 236 209", // #f9ecd1
		"--color-secondary-300": "245 224 181", // #f5e0b5
		"--color-secondary-400": "238 200 125", // #eec87d
		"--color-secondary-500": "230 177 69", // #e6b145
		"--color-secondary-600": "207 159 62", // #cf9f3e
		"--color-secondary-700": "173 133 52", // #ad8534
		"--color-secondary-800": "138 106 41", // #8a6a29
		"--color-secondary-900": "113 87 34", // #715722
		// tertiary | #9f3450 
		"--color-tertiary-50": "241 225 229", // #f1e1e5
		"--color-tertiary-100": "236 214 220", // #ecd6dc
		"--color-tertiary-200": "231 204 211", // #e7ccd3
		"--color-tertiary-300": "217 174 185", // #d9aeb9
		"--color-tertiary-400": "188 113 133", // #bc7185
		"--color-tertiary-500": "159 52 80", // #9f3450
		"--color-tertiary-600": "143 47 72", // #8f2f48
		"--color-tertiary-700": "119 39 60", // #77273c
		"--color-tertiary-800": "95 31 48", // #5f1f30
		"--color-tertiary-900": "78 25 39", // #4e1927
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
		// warning | #EAB308 
		"--color-warning-50": "252 244 218", // #fcf4da
		"--color-warning-100": "251 240 206", // #fbf0ce
		"--color-warning-200": "250 236 193", // #faecc1
		"--color-warning-300": "247 225 156", // #f7e19c
		"--color-warning-400": "240 202 82", // #f0ca52
		"--color-warning-500": "234 179 8", // #EAB308
		"--color-warning-600": "211 161 7", // #d3a107
		"--color-warning-700": "176 134 6", // #b08606
		"--color-warning-800": "140 107 5", // #8c6b05
		"--color-warning-900": "115 88 4", // #735804
		// error | #D41976 
		"--color-error-50": "249 221 234", // #f9ddea
		"--color-error-100": "246 209 228", // #f6d1e4
		"--color-error-200": "244 198 221", // #f4c6dd
		"--color-error-300": "238 163 200", // #eea3c8
		"--color-error-400": "225 94 159", // #e15e9f
		"--color-error-500": "212 25 118", // #D41976
		"--color-error-600": "191 23 106", // #bf176a
		"--color-error-700": "159 19 89", // #9f1359
		"--color-error-800": "127 15 71", // #7f0f47
		"--color-error-900": "104 12 58", // #680c3a
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