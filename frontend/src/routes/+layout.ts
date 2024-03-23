import { locale, loadTranslations } from '$lib/translations/translations';
import type { LayoutLoad } from './$types';

export const ssr = false;
export const prerender = "auto";

export const load: LayoutLoad = async ({ url }: { url: URL }) => {
  const { pathname } = url;
  const defaultLocale = 'en'; // get from cookie / user session etc...
  const initLocale: string = locale.get() || defaultLocale;

  await loadTranslations(initLocale, pathname); // keep this just before the `return`

  return {}
};

