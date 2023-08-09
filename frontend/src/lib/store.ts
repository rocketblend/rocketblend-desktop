import { writable  } from 'svelte/store';
import type { Writable } from 'svelte/store';
import type { project } from './wailsjs/go/models';

export const selectedProject: Writable<project.Project | null> = writable(null)