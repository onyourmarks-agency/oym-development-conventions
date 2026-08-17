---
name: oym-svelte
description: OYM frontend conventions for Svelte 5. Use when writing Svelte components, frontend TypeScript, or Tailwind in On Your Marks projects - both SvelteKit apps and Vite islands embedded in Craft CMS. Covers runes-only state, typing, per-family formatting, folder structure, and Vitest policy.
---

# OYM Svelte

## Step 1: detect the family

One check: `@sveltejs/kit` in `package.json`.

- Present: SvelteKit app. Read `${CLAUDE_PLUGIN_ROOT}/conventions/svelte-universal.md` + `conventions/sveltekit.md`.
- Absent (Vite + Svelte inside a Craft repo, code in `resources/js`): Vite islands. Read `conventions/svelte-universal.md` + `conventions/svelte-vite-craft.md`.

## Warning: the formatting split is intentional

Vite family: semicolons, 2-space indent, trailingComma all. SvelteKit family: tabs, no semicolons, trailingComma none. Both single quotes, width 100. The repo's own Prettier config is authoritative; run the project's format script and never restyle toward the other family.

## Mode

`AGENTS.md` block header `mode=` wins; default `existing` for repos with substantial code. In existing mode mirror the local component tree and idioms (`nevobo-trainersplatform-fe` predates the standard: no Tailwind, BEM CSS; its local style wins inside it).

## Non-negotiables (both families)

- Svelte 5 runes only: `$state`/`$derived`/`$effect`/`$props`. No `$:`, no `export let`, no `on:click` (use `onclick`), no `svelte/store`.
- `<script lang="ts">` always; TS strict; no `any`; props typed inline in the `$props()` destructuring.
- Shared state in `[name]State.svelte.ts` runes files.
- Arrow functions assigned to `const`; explicit braces; full variable names.
- Tailwind 4 CSS-first (`@theme` in CSS); check the existing component tree (shadcn `ui/` or atoms) before writing a new primitive.
- API calls via the project's service classes (native fetch, throw on `!response.ok`); no ad-hoc fetch in components; no axios.
- All UI strings through the project's i18n (i18next or Paraglide).
- SvelteKit: strict client/server boundary; secrets and `$env/dynamic/private` only in `lib/server`/`*.server.ts`; services throw, components catch.
- Run the project's gate before finishing: `pnpm frontend-cs` where it exists, otherwise the lint + format + check scripts (`ddev pnpm ...` in DDEV projects). Never `--no-verify`.
- Vitest unit tests only (entities, states, pure utils); no component-render tests. Add specs alongside new logic; do not scaffold suites in existing repos unprompted.
- Comments: non-obvious "why" only; default AI comment volume is too high.

## New apps

SvelteKit: clone the `oym-chant-fe` structure. Vite islands: scaffold from `oym-vite-boilerplate`. Wire Vitest from `templates/svelte/`. Full checklist: the `oym-new-project` skill.
