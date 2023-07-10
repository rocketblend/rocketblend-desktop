<script>
    import { t } from '$lib/translations/translations';
    import { Greet } from "$lib/wailsjs/go/application/Driver.js";

    let link = 'https://kit.svelte.dev';

    let resultText = "Please enter your name";
    let name = "";

    function greet() {
        Greet(name).then((result) => (resultText = result));
    }
</script>

<main class="container mx-auto p-8 space-y-8">
    <h2>{$t('home.title')}</h2>
    <p>{@html $t('home.text', { link })}</p>
    <hr />
    <section class="card p-4">
        <div class="input-group input-group-divider grid-cols-[auto_1fr_auto]">
            <div class="input-group-shim">{resultText}</div>
            <input
                autocomplete="off"
                bind:value={name}
                class="input"
                id="name"
                type="text"
            />
            <button class="btn variant-filled-primary btn-base" on:click={greet}
                >Greet</button
            >
        </div>
    </section>
    <hr />
    <section class="flex space-x-2">
        <a
            class="btn variant-filled-primary"
            href="https://kit.svelte.dev/"
            target="_blank"
            rel="noreferrer">SvelteKit</a
        >
        <a
            class="btn variant-filled-secondary"
            href="https://tailwindcss.com/"
            target="_blank"
            rel="noreferrer">Tailwind</a
        >
        <a
            class="btn variant-filled-tertiary"
            href="https://github.com/"
            target="_blank"
            rel="noreferrer">GitHub</a
        >
    </section>
</main>
