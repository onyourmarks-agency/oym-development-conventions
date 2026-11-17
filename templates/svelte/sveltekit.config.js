// OYM SvelteKit options. SvelteKit takes its configuration through the Vite
// plugin, so it is here once for vite.config.ts, vitest.config.ts and ESLint:
//   sveltekit(sveltekitConfig)
import adapter from '@sveltejs/adapter-node';

/** @type {import('@sveltejs/kit/vite').Config} */
export const sveltekitConfig = {
	// Runes everywhere in our own code; libraries decide for themselves.
	dynamicCompileOptions: ({ filename }) =>
		filename.includes('node_modules') ? undefined : { runes: true },
	adapter: adapter(),
	csp: {
		directives: {
			'default-src': ['self'],
			'script-src': ['self'],
			'style-src': ['self', 'unsafe-inline'],
			'img-src': ['self', 'data:'],
			'connect-src': ['self'],
			'frame-ancestors': ['none'],
			'base-uri': ['self'],
			'form-action': ['self']
		}
	}
};
