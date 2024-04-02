<script lang="ts">
    import { t } from '$lib/translations/translations';
</script>

<main class="flex flex-col h-full space-y-8">
    <div class="flex flex-col space-y-4">
        <h2 class="h2 font-bold">Welcome!</h2>
        <hr>
    </div>
    <!-- <div class="flex flex-col gap-2">
        <h3 class="text-sm text-surface-800-100-token">Recent projects</h3>
        <div class="grid gap-2 grid-cols-2 lg:grid-cols-3">
            <a class="flex items-center justify-center card p-4" href="/projects">
                <div class="p-4">
                    <h5 class="h5 font-medium">Project 1</h5>
                </div>
            </a>
        </div>
    </div> -->
    <div class="flex flex-col gap-2">
        <h3 class="text-sm text-surface-800-100-token">Shortcuts</h3>
        <div class="grid gap-2 grid-cols-2 lg:grid-cols-4">
            <a class="flex items-center justify-center card p-4" href="/projects">
                <div class="p-4">
                    <h5 class="h5 font-medium">Projects</h5>
                </div>
            </a>
            <a class="flex items-center justify-center card p-4" href="/preferences">
                <div class="p-4">
                    <h5 class="h5 font-medium">Preferences</h5>
                </div>
            </a>
            <!-- <a class="flex items-center justify-center card p-4" href="/getting-started">
                <div class="p-4">
                    <h5 class="h5 font-medium">Getting Started</h5>
                </div>
            </a>
            <a class="flex items-center justify-center card p-4" href="/about">
                <div class="p-4">
                    <h5 class="h5 font-medium">About</h5>
                </div>
            </a>
            <a class="flex items-center justify-center card p-4" href="/projects">
                <div class="p-4">
                    <h5 class="h5 font-medium">New Project</h5>
                </div>
            </a>
        </div> -->
    </div>
</main>