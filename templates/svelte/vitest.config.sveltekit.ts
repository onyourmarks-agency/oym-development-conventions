// OYM Vitest config, SvelteKit family.
// Unit tests only: entities, states, server modules, pure utils under src/lib.
// No component-render tests, no browser environment.
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

import { sveltekitConfig } from './sveltekit.config.js';

export default defineConfig({
	plugins: [sveltekit(sveltekitConfig)],
	test: {
		include: ['src/**/*.spec.ts'],
		environment: 'node',
		globals: false
	}
});
