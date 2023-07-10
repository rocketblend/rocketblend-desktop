import { sveltekit } from '@sveltejs/kit/vite';
import Icons from 'unplugin-icons/vite'

// import path from 'path'

/** @type {import('vite').UserConfig} */
const config = {
  server: {
    fs: {
      // To allow serving files from the frontend project root.
      //
      // allow: ['.'],
    },
  },
  plugins: [
    sveltekit(),
    Icons({
      compiler: 'svelte',
      autoInstall: true,
    }),
  ],
  test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	},
  resolve: {
    alias: {
      'node-fetch': 'isomorphic-fetch',
      // This alias finishes the ability to reference our
      // frontend dirctory with "@path/to/file."
      // You also need to add the path to jsconfig.json.
      //
      // '@': path.resolve(__dirname, './'), 
    },
  },
};

export default config;

