const config = {
  content: [
    "./src/**/*.{html,js,svelte,ts}",
    require('path').join(require.resolve('@skeletonlabs/skeleton'), '../**/*.{html,js,svelte,ts}')
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
  darkMode: 'class',

  theme: {
    extend: {},
  },

  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    ...require('@skeletonlabs/skeleton/tailwind/skeleton.cjs')()
  ],
};

module.exports = config;
