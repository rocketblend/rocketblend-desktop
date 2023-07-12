<script lang="ts">
    import { Table, tableMapperValues } from '@skeletonlabs/skeleton';
    import type { TableSource } from '@skeletonlabs/skeleton';
	import ProjectGallery from '../project-gallery/project-gallery.svelte';

    const sourceData = [
        { position: 1, name: 'Hydrogen', build: "3.4.1", addons: 4, lastOpened: "30 mins ago", bookmarked: true},
        { position: 2, name: 'Helium', build: "3.1.4", addons: 7, lastOpened: "2 mins ago", bookmarked: false},
        { position: 3, name: 'Lithium', build: "3.5.4", addons: 14, lastOpened: "7 weeks ago", bookmarked: false},
        { position: 4, name: 'Beryllium', build: "2.9.12", addons: 0, lastOpened: "2 hours ago", bookmarked: false},
        { position: 5, name: 'Boron', build: "3.4.1", addons: 2, lastOpened: "5 days ago", bookmarked: false},
    ];

    const tableSimple: TableSource = {
        // A list of heading labels.
        head: ['Project', 'Build', 'Addons' ,'Bookmarked', 'Last Opened'],
        // The data visibly shown in your table body UI.
        body: tableMapperValues(sourceData, ['name', 'build', 'addons', 'bookmarked', 'lastOpened']),
        // Optional: The data returned when interactive is enabled and a row is clicked.
        meta: tableMapperValues(sourceData, ['name', 'build']),
    };
</script>

<!-- <Table source={tableSimple} interactive={true} /> -->
<ProjectGallery/>
