<svelte:head>
    <script>
        if (!window.hasOwnProperty("wailsbindings")) {
            let wails_ipc = document.createElement("script");
            wails_ipc.setAttribute("src", "/wails/ipc.js");

            let wails_runtime = document.createElement("script");
            wails_runtime.setAttribute("src", "/wails/runtime.js");

            document.head.appendChild(wails_ipc);
            document.head.appendChild(wails_runtime);
        }
    </script>
</svelte:head>

<slot />
