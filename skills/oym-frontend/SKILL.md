---
name: oym-frontend
description: OYM frontend conventions. Use when writing Svelte components, frontend TypeScript, Tailwind, or frontend build config in On Your Marks projects - SvelteKit apps, Vite islands inside Craft CMS, or legacy webpack islands. Covers runes-only state, typing, prop and Tailwind token rules, per-tool structure, and Vitest policy.
---

# OYM frontend

Layered: read `${CLAUDE_PLUGIN_ROOT}/conventions/shared.md` + `conventions/javascript.md` + `conventions/framework-svelte.md` + exactly ONE tool file.

## Step 1: detect the tool

| Signal | Read |
|---|---|
| `@sveltejs/kit` in package.json | `conventions/tool-sveltekit.md` |
| `webpack/` directory at project root (source under `assets/`) | `conventions/tool-webpack.md` |
| Otherwise (Vite + Svelte in a Craft repo, source under `resources/`) | `conventions/tool-vite.md` |

## Step 2: detect the mode

`AGENTS.md` block header `mode=` wins; default `existing` for repos with substantial code. In existing mode mirror the local component tree and idioms; some older apps predate Tailwind (BEM CSS) and their local style wins inside them. Webpack projects are always `existing`; new projects use Vite or SvelteKit.

## Warning: formatting differs per tool family on purpose

Vite islands: semicolons, 2-space, trailingComma all. SvelteKit: tabs, trailingComma none (semicolons stay on). Both single quotes, width 100. The repo's own Prettier config is authoritative; run the project's format script and never restyle toward another family.

## Non-negotiables

- Runes exclusively: `$state`/`$derived`/`$effect`/`$props`. No `$:`, no `export let`, no `on:click` (use `onclick`), no `svelte/store`.
- `<script lang="ts">` always; TS strict; no `any`.
- Props typed inline in the `$props()` destructuring: each prop on its own line, optional props after required with `?:` and a default.
- Shared state in `[name]State.svelte.ts` runes files.
- Functions as `const fn = () => {}`; braces on every statement body (concise arrow bodies exempt); full variable names.
- Tailwind token projects: no arbitrary bracket values, classes derive from the `@theme` tokens (spacing = px/4), hover via `hover:`/`group-hover:` never JS, missing values become tokens or nearest + uppercase TODO. Never copy design-tool (Figma) output verbatim.
- Check the existing primitives (`ui/` or `atoms/`) before writing a new component; all UI strings through the project's i18n.
- API calls via the project's service classes; services throw, components catch; no ad-hoc fetch in components.
- Run the project's gate before finishing: `pnpm frontend-cs` where it exists, otherwise lint + check (`ddev pnpm ...` in DDEV projects). Never `--no-verify`.
- Vitest unit tests only (entities, states, pure utils); no component-render tests. Add specs alongside new logic; do not scaffold suites in existing repos unprompted.
- Comments: non-obvious "why" only; default AI comment volume is too high.
- Versions come from the repo's manifests, org targets from `${CLAUDE_PLUGIN_ROOT}/VERSIONS.md`; never assert versions from memory.

## New apps

Vite islands scaffold from the `oym-vite-boilerplate` repo; SvelteKit apps follow the structure in `tool-sveltekit.md` with Vitest from `templates/svelte/`. Full checklist: the `oym-new-project` skill.
