<script lang="ts">
    import { t } from '$lib/translations/translations';

    import IconErrorWarningFill from '~icons/ri/error-warning-fill'
</script>

<aside class="alert variant-soft-error">
    <IconErrorWarningFill class="text-4xl"/>
    <div class="alert-message ">
        <h2 class="font-bold h6">Project Directory Not Set</h2>
        <p class="text-sm">Rocketblend Desktop requires a set project directory to monitor changes and create new projects into. Please specify this directory in the preferences page.</p>
    </div>
    <div class="alert-actions">
        <a class="btn btn-sm variant-ghost-error font-medium" href="/preferences#project-watch-directory">Update in Preferences</a>
    </div>
</aside>