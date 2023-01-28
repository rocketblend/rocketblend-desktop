import { c as create_ssr_component, e as escape, d as add_attribute } from "../../chunks/index.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let resultText = "Please enter your name";
  let name = "";
  return `<main class="${"container mx-auto p-8 space-y-8"}"><h1>Welcome to the Unofficial Wails.io SvelteKit Template!</h1>
    <p>Visit <a href="${"https://kit.svelte.dev"}">kit.svelte.dev</a>
        to read the documentation
    </p>
    <hr>
    <section class="${"card p-4"}"><div class="${"input-group input-group-divider grid-cols-[auto_1fr_auto]"}"><div class="${"input-group-shim"}">${escape(resultText)}</div>
            <input autocomplete="${"off"}" class="${"input"}" id="${"name"}" type="${"text"}"${add_attribute("value", name, 0)}>
            <button class="${"btn variant-filled-primary btn-base"}">Greet</button></div></section>
    <hr>
    <section class="${"flex space-x-2"}"><a class="${"btn variant-filled-primary"}" href="${"https://kit.svelte.dev/"}" target="${"_blank"}" rel="${"noreferrer"}">SvelteKit</a>
        <a class="${"btn variant-filled-secondary"}" href="${"https://tailwindcss.com/"}" target="${"_blank"}" rel="${"noreferrer"}">Tailwind</a>
        <a class="${"btn variant-filled-tertiary"}" href="${"https://github.com/"}" target="${"_blank"}" rel="${"noreferrer"}">GitHub</a></section></main>`;
});
export {
  Page as default
};
