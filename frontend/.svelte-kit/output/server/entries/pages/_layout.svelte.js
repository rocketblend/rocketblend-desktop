import { c as create_ssr_component } from "../../chunks/index.js";
const theme = "";
const all = "";
const app = "";
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${$$result.head += `<!-- HEAD_svelte-1rfeo12_START --><script>if (!window.hasOwnProperty("wailsbindings")) {
      let wails_ipc = document.createElement("script");
      wails_ipc.setAttribute("src", "/wails/ipc.js");

      let wails_runtime = document.createElement("script");
      wails_runtime.setAttribute("src", "/wails/runtime.js");

      document.head.appendChild(wails_ipc);
      document.head.appendChild(wails_runtime);
    }
  <\/script><!-- HEAD_svelte-1rfeo12_END -->`, ""}

${slots.default ? slots.default({}) : ``}`;
});
export {
  Layout as default
};
