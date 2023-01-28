<script>
    import { Greet } from "$lib/wailsjs/go/main/App.js";

    let resultText = "Please enter your name below 👇";
    let name = "";

    function greet() {
        Greet(name).then((result) => (resultText = result));
    }
</script>

<main>
    <h1>Welcome to the Unofficial Wails.io SvelteKit Template!</h1>
    <p>
        Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a>
        to read the documentation
    </p>
    <div class="result" id="result">{resultText}</div>
    <div class="input-box" id="input">
        <input
            autocomplete="off"
            bind:value={name}
            class="input"
            id="name"
            type="text"
        />
        <button class="btn" on:click={greet}>Greet</button>
    </div>
</main>

<style>
    .result {
        height: 20px;
        line-height: 20px;
        margin: 1.5rem auto;
    }

    .input-box .btn {
        width: 60px;
        height: 30px;
        line-height: 30px;
        border-radius: 3px;
        border: none;
        margin: 0 0 0 20px;
        padding: 0 8px;
        cursor: pointer;
    }

    .input-box .btn:hover {
        background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
        color: #333333;
    }

    .input-box .input {
        border: none;
        border-radius: 3px;
        outline: none;
        height: 30px;
        line-height: 30px;
        padding: 0 10px;
        background-color: rgba(240, 240, 240, 1);
        -webkit-font-smoothing: antialiased;
    }

    .input-box .input:hover {
        border: none;
        background-color: rgba(255, 255, 255, 1);
    }

    .input-box .input:focus {
        border: none;
        background-color: rgba(255, 255, 255, 1);
    }
</style>
