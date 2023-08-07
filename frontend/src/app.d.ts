/// <reference types="@sveltejs/kit" />
/// <reference types="unplugin-icons/types/svelte" />

declare namespace App {
	// Squelch warnings of image imports from your assets dir
	declare module '$lib/assets/*' {
		const meta: Object[]
		export default meta
	}
	// interface Error {}
	// interface Locals {}
	// interface PageData {}
	// interface Platform {}
}