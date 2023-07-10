import i18n from 'sveltekit-i18n';
import lang from './lang.json';
import type { Config } from 'sveltekit-i18n';
import { load } from '../../../.svelte-kit/types/src/routes/proxy+layout';

interface Params {
  link: string;
  // add more parameters that are used here
}

const config: Config<Params> = {
  translations: {
    en: { lang },
  },
  loaders: [
    {
      locale: 'en',
      key: 'site',
      loader: async () => (
        await import('./en/site.json')
      ).default,
    },
    {
      locale: 'en',
      key: 'home',
      routes: ['/'],
      loader: async () => (
        await import('./en/home.json')
      ).default,
    },
  ],
};

export const { t, locale, locales, loading, loadTranslations } = new i18n(config);
loading.subscribe(($loading) => $loading && console.log('Loading translations...'));