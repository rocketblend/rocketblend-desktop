import { writable  } from 'svelte/store';
import type { Writable } from 'svelte/store';
import type { projectservice } from './wailsjs/go/models';

export const selectedProject: Writable<projectservice.Project | null> = writable(null)