<svelte:head>
    <!-- <script>
        if (!window.hasOwnProperty("wailsbindings")) {
            let wails_app = document.createElement("script");
            wails_app.setAttribute("src", "/wailsjs/ipc.js");

            let wails_runtime = document.createElement("script");
            wails_runtime.setAttribute("src", "/wailsjs/runtime.js");

            document.head.appendChild(wails_app);
            document.head.appendChild(wails_runtime);
        }
    </script> -->
</svelte:head>

<slot />
