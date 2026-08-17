# Svelte conventions (both families)

Applies to all frontend code. OYM has two frontend families; detect first, then read the family file:

- `package.json` depends on `@sveltejs/kit`: SvelteKit app, read `sveltekit.md`.
- Otherwise (Vite + Svelte inside a Craft repo, `resources/` directory): Vite islands, read `svelte-vite-craft.md`.

## Non-negotiables

- `[always]` Svelte 5 runes only: `$state`, `$derived`, `$effect`, `$props`. No legacy `$:` reactive statements, no `export let`, no `on:click` (use `onclick`), and no `svelte/store` (`writable`/`readable`/`derived`) in new code.
- `[always]` Every component is `<script lang="ts">`. TypeScript `strict: true` plus `checkJs`; never disable.
- `[always]` No `any`; use `unknown`, generics, or proper types. Non-null assertions need a reason comment.
- `[always]` Props typed inline in the destructuring: `const { data, size = 'regular' }: { data: Team; size?: string } = $props()`.
- `[always]` Shared reactive state lives in `[name]State.svelte.ts` files using runes (`export const FooState = $state<...>({...})` or a runes class).
- `[always]` Arrow functions assigned to `const`; explicit braces always; full descriptive variable names.
- `[always]` Tailwind 4 CSS-first: theme via `@theme` in the main CSS file; no `tailwind.config.js`.
- `[always]` pnpm, Node 24, ESM. Flat ESLint config. Husky + lint-staged pre-commit; never `--no-verify`.
- `[always]` API calls through the project's service classes/fetch wrappers (native `fetch`, throw on `!response.ok`); no axios, no ad-hoc fetches in components.
- `[always]` No hardcoded UI strings where the project has i18n (i18next in the Vite family, Paraglide in SvelteKit).

## Formatting: intentionally split per family

The two families use opposite Prettier styles. This is deliberate; do not unify, do not "fix".

| | Vite-in-Craft | SvelteKit |
|---|---|---|
| Semicolons | yes | no |
| Indent | 2 spaces | tabs |
| trailingComma | `all` | `none` |
| printWidth | 100 | 100 |
| singleQuote | yes | yes |

- `[always]` The repo's own Prettier/ESLint config is authoritative. Run the project's format and lint scripts; never restyle by hand or import the other family's config.

## Testing

- `[always]` Vitest, unit tests only: entities, state classes, pure utils, service logic with fetch mocked. Never component-rendering or browser tests.
- `[new]` Wire Vitest from day one (`templates/svelte/`).
- `[existing]` Frontend tests are essentially absent today. Same policy as PHP: add specs alongside new logic-bearing code; do not scaffold retroactive suites unprompted.

<!-- oym-card:begin stack=svelte-universal -->
# OYM Svelte rules (both families)

- Svelte 5 runes only: `$state`/`$derived`/`$effect`/`$props`. No `$:`, no `export let`, no `on:click` (use `onclick`), no `svelte/store`.
- `<script lang="ts">` always; TS `strict: true`; no `any` (use `unknown`/generics); non-null assertions need a reason comment.
- Props typed inline in the `$props()` destructuring.
- Shared state in `[name]State.svelte.ts` runes files.
- Arrow functions assigned to `const`; explicit braces; full variable names (no `o`, `a`, `n`).
- Tailwind 4 CSS-first (`@theme` in CSS, no tailwind.config.js).
- pnpm, Node 24, ESM, flat ESLint. Pre-commit hooks are mandatory; never `--no-verify`.
- API calls via the project's service classes (native fetch, throw on `!response.ok`); never ad-hoc fetch in components; no axios.
- Use the project's i18n (i18next or Paraglide); no hardcoded UI strings.
- Formatting differs per family ON PURPOSE (Vite family: semicolons + 2-space; SvelteKit: tabs + no semicolons). The repo's Prettier config is right; run it, never restyle.
- Vitest unit tests only (entities, states, utils); no component-render tests. Existing projects: add specs with new logic only.
<!-- oym-card:end stack=svelte-universal -->

## References

- Vite family canonical: `oym-vite-boilerplate`
- SvelteKit canonical: `oym-chant-fe`
