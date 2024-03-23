<script lang="ts">
    import { dev } from '$app/environment';
    
    import SuperDebug from "sveltekit-superforms";
    import { type SuperValidated, type Infer, superForm, setMessage, setError } from "sveltekit-superforms";
    import { zodClient } from "sveltekit-superforms/adapters";

    import { ProgressRadial } from '@skeletonlabs/skeleton'; 
    import { desktopGeneralForm, type DesktopGeneralForm } from "./schema";

    export let data: SuperValidated<Infer<DesktopGeneralForm>>;
 
    const form = superForm(data, {
        validators: zodClient(desktopGeneralForm),
        delayMs: 1250,
        SPA: true,
        onUpdate({ form }) {
            if (form.valid) {
                setMessage(form, "Your preferences have been saved.");
            }
        },
    });

    const { form: formData, errors, message, constraints, enhance, delayed } = form;
</script>

<div>
    {#if $message}
        <aside class="alert variant-filled-success max-w-md">
            <div class="alert-message">
                <h3 class="h3">Thank you!</h3>
                <p>{$message}</p>
            </div>
        </aside>
    {:else}
        <form method="POST" class="space-y-2" use:enhance>
            <div>
                <h5 class="h5 font-bold">Projects</h5>
                <label class="label">
                    <p class="text-sm text-surface-200">Project directory - Folder used to populate projects and watch for changes.</p>
                    <input
                        class="input"
                        type="text"
                        aria-invalid={$errors.watchFolder ? 'true' : undefined}
                        bind:value={$formData.watchFolder}
                        {...$constraints.watchFolder}
                    />
                    {#if $errors.watchFolder}<span class="invalid">{$errors.watchFolder}</span>{/if}
                </label>
            </div>

            <div>
                <h5 class="h5 font-bold">Features</h5>
            </div>

            <button type="button" disabled={$delayed} class="btn btn-sm variant-filled">
                {#if $delayed}
                    <ProgressRadial class="mr-2 h-4 w-4"/>
                    Saving...
                {:else}
                    Save
                {/if}
            </button>
        
            {#if dev}
                <SuperDebug data={$formData} />
            {/if}
        </form>
    {/if}
</div>