import { c as create_ssr_component, e as escape, d as add_attribute } from "../../chunks/index.js";
const _page_svelte_svelte_type_style_lang = "";
const css = {
  code: ".result.svelte-19yygje.svelte-19yygje{height:20px;line-height:20px;margin:1.5rem auto}.input-box.svelte-19yygje .btn.svelte-19yygje{width:60px;height:30px;line-height:30px;border-radius:3px;border:none;margin:0 0 0 20px;padding:0 8px;cursor:pointer}.input-box.svelte-19yygje .btn.svelte-19yygje:hover{background-image:linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);color:#333333}.input-box.svelte-19yygje .input.svelte-19yygje{border:none;border-radius:3px;outline:none;height:30px;line-height:30px;padding:0 10px;background-color:rgba(240, 240, 240, 1);-webkit-font-smoothing:antialiased}.input-box.svelte-19yygje .input.svelte-19yygje:hover{border:none;background-color:rgba(255, 255, 255, 1)}.input-box.svelte-19yygje .input.svelte-19yygje:focus{border:none;background-color:rgba(255, 255, 255, 1)}",
  map: null
};
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let resultText = "Please enter your name below 👇";
  let name = "";
  $$result.css.add(css);
  return `<main><h1>Welcome to the Unofficial Wails.io SvelteKit Template!</h1>
    <p>Visit <a href="${"https://kit.svelte.dev"}">kit.svelte.dev</a>
        to read the documentation
    </p>
    <div class="${"result svelte-19yygje"}" id="${"result"}">${escape(resultText)}</div>
    <div class="${"input-box svelte-19yygje"}" id="${"input"}"><input autocomplete="${"off"}" class="${"input svelte-19yygje"}" id="${"name"}" type="${"text"}"${add_attribute("value", name, 0)}>
        <button class="${"btn svelte-19yygje"}">Greet</button></div>
</main>`;
});
export {
  Page as default
};
