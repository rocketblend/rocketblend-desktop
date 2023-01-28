import { c as create_ssr_component } from "../../chunks/index.js";
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${$$result.head += `<!-- HEAD_svelte-nitc86_START --><script>if (!window.hasOwnProperty("wailsbindings")) {
            let wails_ipc = document.createElement("script");
            wails_ipc.setAttribute("src", "/wails/ipc.js");

            let wails_runtime = document.createElement("script");
            wails_runtime.setAttribute("src", "/wails/runtime.js");

            document.head.appendChild(wails_ipc);
            document.head.appendChild(wails_runtime);
        }
    <\/script><!-- HEAD_svelte-nitc86_END -->`, ""}

${slots.default ? slots.default({}) : ``}`;
});
export {
  Layout as default
};
