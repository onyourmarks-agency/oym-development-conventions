# Webpack islands conventions (tool layer)

Applies to older Craft projects whose frontend is built with the webpack boilerplate lineage. Builds on `javascript.md` + `framework-svelte.md`. This tool is `[existing]` only: no new projects start on webpack; new island frontends use `tool-vite.md`.

## Recognizing a webpack project

A `webpack/` directory at the project root (`webpack.config.js`, `config/`, `plugins/`, `pre-script.js`) and frontend source under `assets/` instead of `resources/`.

## Rules

- `[existing]` Source lives under `assets/` (js, css/scss, svg). Respect the existing folder layout; these projects predate the atomic-design tree.
- `[existing]` Build scripts run the pre-script first: `pnpm watch` (dev), `pnpm build` (production). The pre-script generates required files; never bypass it, and run it before lint tasks that depend on generated output (the repo's scripts already chain it).
- `[existing]` The same aggregate gate applies: `pnpm frontend-cs` (format + eslint + stylelint + svelte-check). Stylelint here covers SCSS with BEM patterns.
- `[existing]` Styling in these projects is SCSS/BEM era, not Tailwind; follow the local component and class naming exactly.
- `[existing]` Svelte components still follow `framework-svelte.md` (runes, lang="ts", prop rules) for new code, but match the surrounding file when editing older non-runes components; converting a component to runes is a task, not a drive-by.
- `[existing]` Do not migrate a webpack project to Vite as a side effect of another task. The migration is its own project.

<!-- oym-card:begin stack=tool-webpack -->
# OYM webpack islands rules (legacy Craft frontends)

- Existing projects only; new island frontends use Vite.
- Source under `assets/`, build via the `webpack/` directory; always run through the pre-script chain (`pnpm watch`, `pnpm build`), never bypass it.
- Gate: `pnpm frontend-cs` (format + eslint + stylelint + svelte-check); stylelint covers SCSS/BEM.
- Styling is SCSS/BEM, not Tailwind; follow local class naming exactly.
- New Svelte code follows the Svelte layer (runes, lang="ts"); match surrounding style when editing older components; runes conversions and Vite migrations are dedicated tasks, never drive-bys.
<!-- oym-card:end stack=tool-webpack -->

